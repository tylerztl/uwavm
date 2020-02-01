BUILD_DIR ?= build

DRIVER_SRCS := $(shell ls ./driver/*.cc)
DRIVER_OBJS := $(DRIVER_SRCS:%=$(BUILD_DIR)/%.o)
DRIVER_DEPS := $(DRIVER_OBJS:.o=.d)

CONTRACT_SRCS := $(shell ls ./example/*.cc)
CONTRACT_BIN := $(CONTRACT_SRCS:./example/%.cc=$(BUILD_DIR)/%.wasm)

CONTRACT_PB_SRCS := $(shell ls ./pb/*.cc)
CONTRACT_PB_OBJS := $(CONTRACT_PB_SRCS:./pb/%.cc=$(BUILD_DIR)/pb/%.cc.o)

INC_DIRS := . /usr/local/include
INC_FLAGS := $(addprefix -I,$(INC_DIRS))

OPT_LEVEL := -Oz
ifeq ($(CXX),g++)
	OPT_LEVEL := -Os
endif

CPPFLAGS ?= $(INC_FLAGS) -MMD -MP $(OPT_LEVEL) -std=c++11

.PHONY: clean all

all: $(BUILD_DIR)/libdriver.a $(CONTRACT_BIN)

clean:
	$(RM) -r $(BUILD_DIR)

$(BUILD_DIR)/libdriver.a: $(DRIVER_OBJS)
	$(AR) -rc $@ $^
	$(RANLIB) $@

$(BUILD_DIR)/test/%.out: $(BUILD_DIR)/test/%.cc.o $(CONTRACT_PB_OBJS) $(BUILD_DIR)/libdriver.a
	$(CXX) $(CPPFLAGS) $(CXXFLAGS) -o $@ $^ -O3 -L/usr/local/lib -lprotobuf-lite -lpthread -lgtest
	$@

# wasm target
$(BUILD_DIR)/%.wasm: $(BUILD_DIR)/example/%.cc.o $(CONTRACT_PB_OBJS) $(BUILD_DIR)/libdriver.a
	$(CXX) $(CPPFLAGS) $(CXXFLAGS) -o $@ $^ -Oz \
								-s DETERMINISTIC=1 \
								-s TOTAL_STACK=256KB \
								-s TOTAL_MEMORY=1MB \
								-s EXTRA_EXPORTED_RUNTIME_METHODS='["stackAlloc"]' \
								--js-library driver/exports.js \
								-L/usr/local/lib -lprotobuf-lite -lpthread

# c++ source
$(BUILD_DIR)/%.cc.o: %.cc
	$(MKDIR_P) $(dir $@)
	$(CXX) $(CPPFLAGS) $(CXXFLAGS) -c $< -o $@

-include $(DRIVER_DEPS)

MKDIR_P ?= mkdir -p
RANLIB ?= ranlib
