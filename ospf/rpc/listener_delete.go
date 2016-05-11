package rpc

import (
	"fmt"
	"ospfd"
	//    "l3/ospf/config"
	//    "l3/ospf/server"
	//    "utils/logging"
	//    "net"
)

func (h *OSPFHandler) DeleteOspfGlobal(ospfGlobalConf *ospfd.OspfGlobal) (bool, error) {
	h.logger.Info(fmt.Sprintln("Delete global config attrs:", ospfGlobalConf))
	return true, nil
}

func (h *OSPFHandler) DeleteOspfAreaEntry(ospfAreaConf *ospfd.OspfAreaEntry) (bool, error) {
	h.logger.Info(fmt.Sprintln("Delete Area Config attrs:", ospfAreaConf))
	return true, nil
}

func (h *OSPFHandler) DeleteOspfStubAreaEntry(ospfStubAreaConf *ospfd.OspfStubAreaEntry) (bool, error) {
	h.logger.Info(fmt.Sprintln("Delete Stub Area Config attrs:", ospfStubAreaConf))
	return true, nil
}

func (h *OSPFHandler) DeleteOspfIfEntry(ospfIfConf *ospfd.OspfIfEntry) (bool, error) {
	h.logger.Info(fmt.Sprintln("Delete interface config attrs:", ospfIfConf))
	return true, nil
}

func (h *OSPFHandler) DeleteOspfIfMetricEntry(ospfIfMetricConf *ospfd.OspfIfMetricEntry) (bool, error) {
	h.logger.Info(fmt.Sprintln("Delete interface metric config attrs:", ospfIfMetricConf))
	return true, nil
}

func (h *OSPFHandler) DeleteOspfVirtIfEntry(ospfVirtIfConf *ospfd.OspfVirtIfEntry) (bool, error) {
	h.logger.Info(fmt.Sprintln("Delete virtual interface config attrs:", ospfVirtIfConf))
	return true, nil
}
