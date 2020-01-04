## Basic info
When managing big Proxmox VE clusters its hard to remember on which node LXC/KCM instance is located. Of course, you can just use a web-interface to find it, but I think it would be nice to have the ability to do it from the terminal.

## Installation
Just put binary anywhere you need it and give execution permission to it (`chmod +x`). To make it available system-wide you can just put it in someplace from your $PATH (ex. /usr/bin or /usr/local/bin).

***ATTENTION***: You must run pve-search on the cluster node in which you want to perform a search.


## Usage
To search instances across cluster simply run `pve-search [GoRegexp for name]`, e.g. if you want to search cluster for some instance, that have "load-balancer" in its name, just run `pve-search load-balancer`, it will give you something like this:

    +------+------+---------+-----------------+-------------+----------------+------------------------+-------------------------+---------------+---------------+
    | VMID | TYPE |  STATE  |     VM NAME     |  NODE NAME  |      CPU       |      FREE MEMORY       |        FREE DISK        |    NET IN     |    NET OUT    |
    +------+------+---------+-----------------+-------------+----------------+------------------------+-------------------------+---------------+---------------+
    | 8104 | lxc  | running | load-balancer4  | node9       | 4.29% (24 CPU) | 12944MB/16384MB 79.00% | 7108MB/20030MB 35.48%   | 5828797.10 MB | 6431487.52 MB |
    | 8103 | lxc  | running | load-balancer3  | node6       | 4.88% (24 CPU) | 12627MB/16000MB 78.92% | 7050MB/20030MB 35.20%   | 6237860.60 MB | 6855118.29 MB |
    | 8105 | lxc  | running | load-balancer5  | node8       | 6.79% (16 CPU) | 12863MB/16384MB 78.51% | 7040MB/20030MB 35.14%   | 6250774.89 MB | 6868919.56 MB |
    | 8109 | lxc  | stopped | load-balancer9  | node10      | 0.00% (8 CPU)  | 8192MB/8192MB 100.00%  | 20480MB/20480MB 100.00% | 0.00 MB       | 0.00 MB       |
    | 8107 | lxc  | stopped | load-balancer7  | node7       | 0.00% (8 CPU)  | 8192MB/8192MB 100.00%  | 20480MB/20480MB 100.00% | 0.00 MB       | 0.00 MB       |
    | 8110 | lxc  | stopped | load-balancer10 | node2       | 0.00% (8 CPU)  | 8192MB/8192MB 100.00%  | 20480MB/20480MB 100.00% | 0.00 MB       | 0.00 MB       |
    | 8108 | lxc  | stopped | load-balancer8  | node5       | 0.00% (8 CPU)  | 8192MB/8192MB 100.00%  | 20480MB/20480MB 100.00% | 0.00 MB       | 0.00 MB       |
    | 8102 | lxc  | running | load-balancer2  | node3       | 4.30% (24 CPU) | 16557MB/20000MB 82.78% | 7036MB/20030MB 35.13%   | 6243922.93 MB | 6867503.67 MB |
    | 8101 | lxc  | running | load-balancer1  | node1       | 4.89% (24 CPU) | 12976MB/16384MB 79.20% | 6954MB/20030MB 34.72%   | 5757362.84 MB | 6688386.79 MB |
    | 8106 | lxc  | running | load-balancer6  | node4       | 0.03% (8 CPU)  | 6517MB/8192MB 79.55%   | 13272MB/20030MB 66.26%  | 2481.95 MB    | 2145.17 MB    |
    +------+------+---------+-----------------+-------------+----------------+------------------------+-------------------------+---------------+---------------+
 
You can use different sorting keys with `--sort [key]` flag. Valid keys are: cpu, mem, vmname, vmid, node. By default it sorts results in descending order, to sort by ascending use `--asc` flag. Also, you can output values in plain text, using `--asc` flag. You can also apply additional filters with `--state` `--type` and `--vmid` flag followed by goregexp. You can read about regexp syntax [Here](https://github.com/google/re2/wiki/Syntax). 
Also, you can output results in space-separated values using `--text` flag. In this case, values will be in the same order, as in the table, but data values such as free memory or free disk will be in bytes instead of megabytes