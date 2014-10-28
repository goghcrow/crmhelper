package crm2

import "bytes"

func (self *Crm2) GetBuf() *bytes.Buffer {
	var buf *bytes.Buffer
	select {
	case buf = <-self.chBuf:
		buf.Reset()
	default:
		buf = bytes.NewBuffer(make([]byte, 0, self.bufferSize))
	}
	return buf
}

func (self *Crm2) PutBuf(buf *bytes.Buffer) {
	select {
	case self.chBuf <- buf:
	default: // prevent full chBuf block
	}
}
