package data_go

//
//import (
//	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
//	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
//)
//
//func main() {
//	proxywasm.SetNewHttpContext(newContext)
//}
//
//type httpHeaders struct {
//	// you must embed the default context so that you need not to reimplement all the methods by yourself
//	proxywasm.DefaultHttpContext
//	contextID uint32
//}
//
//func newContext(rootContextID, contextID uint32) proxywasm.HttpContext {
//	return &httpHeaders{contextID: contextID}
//}
//
//// override
//func (ctx *httpHeaders) OnHttpRequestHeaders(numHeaders int, endOfStream bool) types.Action {
//	hs, err := proxywasm.GetHttpRequestHeaders()
//	if err != nil {
//		proxywasm.LogCriticalf("failed to get request headers: %v", err)
//	}
//
//	for _, h := range hs {
//		proxywasm.LogInfof("request header --> %s: %s", h[0], h[1])
//	}
//	return types.ActionContinue
//}
//
//// override
//func (ctx *httpHeaders) OnHttpResponseHeaders(numHeaders int, endOfStream bool) types.Action {
//	hs, err := proxywasm.GetHttpResponseHeaders()
//	if err != nil {
//		proxywasm.LogCriticalf("failed to get request headers: %v", err)
//	}
//
//	for _, h := range hs {
//		proxywasm.LogInfof("response header <-- %s: %s", h[0], h[1])
//	}
//	return types.ActionContinue
//}
//
//// override
//func (ctx *httpHeaders) OnHttpStreamDone() {
//	proxywasm.LogInfof("%d finished", ctx.contextID)
//}
