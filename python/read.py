from google.protobuf.internal.encoder import _VarintBytes
from google.protobuf.internal.decoder import _DecodeVarint32

import metric_pb2
import random

with open('out.bin', 'rb') as f:
    buf = f.read()
    n = 0
    count = 0
    l = 0

    #while n < len(buf):
    ptr = buf
    while l < len(buf):
        msg_len, new_pos = _DecodeVarint32(ptr, 0)
        n = new_pos
        msg_buf = ptr[n:n+msg_len]
        n += msg_len
        l += n
        ptr = ptr[n:]
        read_metric = metric_pb2.Metric()
        read_metric.ParseFromString(msg_buf)
        print(read_metric)
        count += 1

    print(count)
