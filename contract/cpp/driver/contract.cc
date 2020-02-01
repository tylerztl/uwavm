#include "driver/context_impl.h"
#include "pb/contract.pb.h"
#include "driver/syscall.h"
#include "driver/driver.h"

namespace uwavm {

static void return_response(const Response& resp) {
    pb::SetOutputRequest req;
    pb::SetOutputResponse rep;
    pb::Response* r = req.mutable_response();
    r->set_status(resp.status);
    r->set_message(resp.message);
    r->set_body(resp.body);
    syscall("SetOutput", req, &rep);
}

Contract::Contract() {
    ContextImpl* ctximpl = new (ContextImpl);
    ctximpl->init();
    _ctx = ctximpl;
}

Contract::~Contract() {
    return_response(*_ctx->mutable_response());
    delete (_ctx);
}

}  // namespace uwavm
