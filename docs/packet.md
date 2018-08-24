请求ip+port   172.16.100.11:7777

#### 包格式
<table>
<tr><th colspan=4  align="center">包头构成</th></tr>
<tr><td>包体长度(int32)</td><td>请求序号(int32)</td><td>命令符(int8)</td><td>服务版本(int8)</td></tr>
<tr><td> 4byte 大端</td><td> 4byte大端</td><td>1byte</td><td>1byte</td></tr>
<tr><td colspan=4 > 10byte</td></tr>
</table>

#### 示例代码
```go
func extractHeader(header []byte) (int, int, int8, int8) {
	var length int32
	var seq int32
	var traceId int32
	buffer := bytes.NewBuffer(header)
	binary.Read(buffer, binary.BigEndian, &length)
	binary.Read(buffer, binary.BigEndian, &seq)
	cmd, _ := buffer.ReadByte()
	version, _ := buffer.ReadByte()
	return int(length), int(seq), int8(cmd), int8(version)
}
```