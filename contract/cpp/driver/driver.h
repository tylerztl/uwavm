#ifndef DRIVER_DRIVER_H
#define DRIVER_DRIVER_H

#include <map>
#include <string>
#include <vector>
#include <memory>

namespace uwavm {

struct Response {
    int status;
    std::string message;
    std::string body;
};

const std::string kUnknownKey = "";

class Context {
public:
    virtual ~Context() {}
    virtual const std::map<std::string, std::string>& args() const = 0;
    virtual const std::string& arg(const std::string& name) const = 0;
    virtual const std::string& caller() const = 0;
    virtual bool get_object(const std::string& key, std::string* value) = 0;
    virtual bool put_object(const std::string& key,
                            const std::string& value) = 0;
    virtual bool delete_object(const std::string& key) = 0;
    virtual void ok(const std::string& body) = 0;
    virtual void error(const std::string& body) = 0;
    virtual Response* mutable_response() = 0;
    virtual bool call(const std::string& module, const std::string& contract,
                      const std::string& method,
                      const std::map<std::string, std::string>& args,
                      Response* response) = 0;
};

class Contract {
public:
    Contract();
    virtual ~Contract();
    Context* context() { return _ctx; };

private:
    Context* _ctx;
};

}  // namespace uwavm

#define DEFINE_METHOD(contract_class, method_name)        \
    static void cxx_##method_name(contract_class&);       \
    extern "C" void __attribute__((used)) method_name() { \
        contract_class self;                              \
        cxx_##method_name(self);                          \
    };                                                    \
    static void cxx_##method_name(contract_class& self)

#endif
