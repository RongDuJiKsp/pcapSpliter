# pcapSpliter

`pcapSpliter` 是一个用于拆分 DoH 抓包文件的工具。它会读取一个离线 pcap 文件，按 HTTPS 流把数据切成多个小 pcap，并把结果写到输入文件同名目录下。

`pcapSpliter` is a tool for splitting DoH capture files. It reads an offline pcap file, breaks the data into multiple smaller pcap files by HTTPS flow, and writes the results into a directory with the same name as the input file.

## Usage

先构建：

Build first:

```bash
go build
```

再传入一个 pcap 文件：

Then pass in a pcap file:

```bash
pcapSpliter <filename>
```

例如：

For example:

```bash
pcapSpliter sample.pcap
```

程序会在 `sample.pcap` 所在目录下创建一个 `sample/` 目录，并把拆分结果输出到里面。

The program creates a `sample/` directory next to `sample.pcap` and writes all split results into that folder.

## Output

输出文件命名规则为：

Output files follow this naming rule:

```text
<原文件名>_<DNS Provider>_<时间戳>.pcap
```

例如：

For example:

```text
sample_Cloudflare_2026-08-03_Mon_12-30-00_CST.pcap
```

如果输入的是 `sample.pcap`，那么最终目录结构大致如下：

If the input is `sample.pcap`, the final directory layout looks like this:

```text
sample/
	sample_Cloudflare_2026-08-03_Mon_12-30-00_CST.pcap
	sample_Google_2026-08-03_Mon_12-31-10_CST.pcap
```

## 底层实现

整体流程是：先离线读取 pcap，再筛出 IPv4 + TCP + 443 的流量，然后按四元组聚合成 flow，最后把每条 DoH 流写成独立的 pcap 文件。

The overall flow is: read the pcap offline, filter for IPv4 + TCP + 443 traffic, group packets into flows by four-tuple, and finally write each DoH flow into its own pcap file.

### 1. 离线读取 pcap

入口在 [main.go](main.go) 中。程序只接受一个参数，也就是待处理的 pcap 文件路径。随后会通过 `gopacket` 以离线模式读取该文件。

The entry point is [main.go](main.go). The program accepts exactly one argument, which is the pcap file to process, and then uses `gopacket` to read it in offline mode.

### 2. 自动创建输出目录

输出目录由 [utils/fileutil/dir.go](utils/fileutil/dir.go) 创建，规则是取输入文件名去掉扩展名，作为同级目录名。例如 `sample.pcap` 会生成 `sample/`。

The output directory is created in [utils/fileutil/dir.go](utils/fileutil/dir.go). It uses the input file name without the extension as the sibling directory name, so `sample.pcap` becomes `sample/`.

### 3. 只处理 IPv4 + TCP + 443 流量

分片逻辑在 [flowmaker/maker.go](flowmaker/maker.go) 和 [flowmaker/assembler.go](flowmaker/assembler.go) 中完成：

- 只处理 IPv4 包
- 只处理 TCP 包
- 只保留源端口或目的端口为 443 的流量
- 通过四元组 `srcIP:srcPort -> dstIP:dstPort` 作为 flow key

The splitting logic lives in [flowmaker/maker.go](flowmaker/maker.go) and [flowmaker/assembler.go](flowmaker/assembler.go):

- IPv4 packets only
- TCP packets only
- Only traffic with source port or destination port equal to 443
- A four-tuple `srcIP:srcPort -> dstIP:dstPort` is used as the flow key

### 4. 识别 DoH 服务商

程序会根据 443 端点的 IP 判断 DNS Provider。当前内置了 Cloudflare、Google、Quad9、Serveroid、AdGuard、Tencent、TWNIC 和 Proxy 等映射，定义在 [utils/dns/nameof.go](utils/dns/nameof.go) 和 [utils/dns/func.go](utils/dns/func.go) 中。

The DNS provider is identified from the IP address on the 443 endpoint. Built-in mappings include Cloudflare, Google, Quad9, Serveroid, AdGuard, Tencent, TWNIC, and Proxy, and they are defined in [utils/dns/nameof.go](utils/dns/nameof.go) and [utils/dns/func.go](utils/dns/func.go).

If the IP is not in the mapping table, the flow is treated as non-DoH traffic and skipped.

### 5. 按流写出新的 pcap 文件

每当程序遇到一个新的 DoH 流，就会创建一个新的 pcap 文件，并把该流里的后续包持续写入该文件。文件名会带上：

- 原始输入文件名
- DNS Provider 名称
- 首个包的时间戳

Whenever the program sees a new DoH flow, it creates a new pcap file and keeps writing subsequent packets from that flow into it. The file name includes:

- The original input file name
- The DNS provider name
- The timestamp of the first packet

此外，代码里还有一个流切分的保护逻辑：当当前流包数已经超过阈值，并且再次看到 SYN 包时，会认为可能进入了新的会话并重新开文件（和目前最有名的数据集捕获工具DoHlyzer的逻辑一致）。

There is also a safeguard for flow splitting: once the current flow grows beyond a threshold and a SYN packet appears again, the tool assumes a new session may have started and opens a new file. This matches the logic used by the best-known dataset tool (DoHlyzer) for this problem.

## Notes

- 目前只支持 IPv4。
- 目前只处理 TCP 443，不处理 UDP/HTTP3。
- 输出目录如果已经存在，`Mkdir` 会直接返回错误，需要手动清理或换一个输入文件名。

- IPv4 only is supported for now.
- Only TCP 443 is handled; UDP/HTTP3 are not processed.
- If the output directory already exists, `Mkdir` returns an error, so you need to clean it up manually or use a different input file name.
