package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
)

type vm struct {
	CPU       float64 `json:"cpu"`
	Disk      uint64  `json:"disk"`
	Diskread  uint64  `json:"diskread"`
	Diskwrite uint64  `json:"diskwrite"`
	ID        string  `json:"id"`
	Maxcpu    int8    `json:"maxcpu"`
	Maxdisk   uint64  `json:"maxdisk"`
	Maxmem    uint64  `json:"maxmem"`
	Mem       uint64  `json:"mem"`
	Name      string  `json:"name"`
	Netin     uint64  `json:"netin"`
	Netout    uint64  `json:"netout"`
	Node      string  `json:"node"`
	Status    string  `json:"status"`
	Template  int8    `json:"template"`
	Type      string  `json:"type"`
	Uptime    int32   `json:"uptime"`
	Vmid      int     `json:"vmid"`
}

func main() {
	vmid := flag.String("vmid", ".*", "add vmid filter")
	state := flag.String("state", ".*", "add state filter")
	vmtype := flag.String("type", ".*", "add vm type filter")
	toSort := flag.String("sort", "nosort", "sort resul by key (cpu, mem, disk, vmname, vmid, node)")
	byAsc := flag.Bool("asc", false, "sort by ascending")
	toText := flag.Bool("text", false, "Output in raw text")
	flag.Parse()
	detectPVEcluster()
	array := getVMarray()
	array = parseArray(array, *vmid, flag.Arg(0), *state, *vmtype)
	switch *toSort {
	case "nosort":
		break
	case "cpu", "mem", "disk", "vmname", "vmid", "node":
		array = sortArray(array, toSort, *byAsc)
	default:
		fmt.Printf("%s is not correct sorting key, use one of these: cpu, mem, disk, vmname, nodename", *toSort)
		os.Exit(10)
	}
	printResult(array, *toText)
}

func detectPVEcluster() {
	_, err := os.Stat("/etc/pve/corosync.conf")
	if err != nil {
		fmt.Println("Corosync config not found in /etc/pve/. Exiting...")
		os.Exit(1)
	}
}

func getVMarray() []vm {
	jsonBulk := make([]vm, 0)
	jsonRaw, err := exec.Command("pvesh", "get", "cluster/resources", "-type", "vm", "-output-format", "json").CombinedOutput()
	check(err)
	json.Unmarshal(jsonRaw, &jsonBulk)
	return jsonBulk
}

func parseArray(arr []vm, vmid string, name string, state string, vmtype string) []vm {
	result := make([]vm, 0)
	if name == "" {
		name = ".*"
	}
	reName, err := regexp.Compile(name)
	checkRegexp(err, name)
	reVmid, err := regexp.Compile(vmid)
	checkRegexp(err, vmid)
	reState, err := regexp.Compile(state)
	checkRegexp(err, state)
	reType, err := regexp.Compile(vmtype)
	checkRegexp(err, vmtype)
	for i := range arr {
		if reName.Match([]byte(arr[i].Name)) && reVmid.Match([]byte(fmt.Sprintf("%d", arr[i].Vmid))) && reState.Match([]byte(arr[i].Status)) && reType.Match([]byte(arr[i].Type)) {
			result = append(result, arr[i])
		}
	}
	return result
}

func sortArray(arr []vm, key *string, byAsc bool) []vm {
	switch *key {
	case "cpu":
		if byAsc {
			sort.SliceStable(arr, func(i, j int) bool { return arr[i].CPU > arr[j].CPU })
		} else {
			sort.SliceStable(arr, func(i, j int) bool { return arr[i].CPU < arr[j].CPU })
		}

	case "mem":
		if byAsc {
			sort.SliceStable(arr, func(i, j int) bool {
				return float64(arr[i].Maxmem-arr[i].Mem)/float64(arr[i].Maxmem) > float64(arr[j].Maxmem-arr[j].Mem)/float64(arr[j].Maxmem)
			})
		} else {
			sort.SliceStable(arr, func(i, j int) bool {
				return float64(arr[i].Maxmem-arr[i].Mem)/float64(arr[i].Maxmem) < float64(arr[j].Maxmem-arr[j].Mem)/float64(arr[j].Maxmem)
			})
		}
	case "disk":
		if byAsc {
			sort.SliceStable(arr, func(i, j int) bool {
				return float64(arr[i].Maxdisk-arr[i].Disk)/float64(arr[i].Maxdisk) > float64(arr[j].Maxdisk-arr[j].Disk)/float64(arr[j].Maxdisk)
			})
		} else {
			sort.SliceStable(arr, func(i, j int) bool {
				return float64(arr[i].Maxdisk-arr[i].Disk)/float64(arr[i].Maxdisk) < float64(arr[j].Maxdisk-arr[j].Disk)/float64(arr[j].Maxdisk)
			})
		}
	case "vmname":
		if byAsc {
			sort.SliceStable(arr, func(i, j int) bool { return arr[i].Name > arr[j].Name })
		} else {
			sort.SliceStable(arr, func(i, j int) bool { return arr[i].Name < arr[j].Name })
		}
	case "vmid":
		if byAsc {
			sort.SliceStable(arr, func(i, j int) bool { return arr[i].Vmid > arr[j].Vmid })
		} else {
			sort.SliceStable(arr, func(i, j int) bool { return arr[i].Vmid < arr[j].Vmid })
		}
	case "node":
		if byAsc {
			sort.SliceStable(arr, func(i, j int) bool { return arr[i].Node > arr[j].Node })
		} else {
			sort.SliceStable(arr, func(i, j int) bool { return arr[i].Node < arr[j].Node })
		}
	default:
		os.Exit(1)
	}
	return arr
}

func printResult(arr []vm, toText bool) {
	if toText {
		for i := range arr {
			cpu := arr[i].CPU * 100
			memoryFree := arr[i].Maxmem - arr[i].Mem
			diskFree := arr[i].Maxdisk - arr[i].Disk
			fmt.Printf("%d %s %s %s %.2f %d %d %d %d %d %d %d %d\n", arr[i].Vmid, arr[i].Type, arr[i].Name, arr[i].Node, cpu, arr[i].Maxmem, arr[i].Mem, memoryFree, arr[i].Maxdisk, arr[i].Disk, diskFree, arr[i].Netin, arr[i].Netout)
		}
	} else {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"vmid", "type", "state", "vm name", "node name", "cpu", "free memory", "free disk", "net in", "net out"})
		for i := range arr {
			cpu := strings.Join([]string{fmt.Sprintf("%.2f%%", arr[i].CPU*100), fmt.Sprintf("(%d CPU)", arr[i].Maxcpu)}, " ")
			memory := strings.Join([]string{fmt.Sprintf("%.0fMB", math.Round(float64(arr[i].Maxmem-arr[i].Mem)/1024/1024)), "/", fmt.Sprintf("%.0fMB", math.Round(float64(arr[i].Maxmem/1024/1024))), fmt.Sprintf(" %.2f%%", float64(arr[i].Maxmem-arr[i].Mem)/float64(arr[i].Maxmem)*100)}, "")
			netin := fmt.Sprintf("%.2f MB", float64(arr[i].Netin)/1024/1024)
			netout := fmt.Sprintf("%.2f MB", float64(arr[i].Netout)/1024/1024)
			diskFree := strings.Join([]string{fmt.Sprintf("%.0fMB", math.Round(float64(arr[i].Maxdisk-arr[i].Disk)/1024/1024)), "/", fmt.Sprintf("%.0fMB", math.Round(float64(arr[i].Maxdisk/1024/1024))), fmt.Sprintf(" %.2f%%", float64(arr[i].Maxdisk-arr[i].Disk)/float64(arr[i].Maxdisk)*100)}, "")
			table.Append([]string{strconv.Itoa(arr[i].Vmid), arr[i].Type, arr[i].Status, arr[i].Name, arr[i].Node, cpu, memory, diskFree, netin, netout})
		}
		table.Render()
	}
}

func checkRegexp(e error, exp string) {
	if e != nil {
		fmt.Println("Whoops...There something wronmg with your regexp! Please check following output. Your regexp was: ", exp)
		panic(e)
	}
}
func check(e error) {
	if e != nil {
		panic(e)
	}
}
