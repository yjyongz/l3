package server

import (
        "encoding/binary"
)

const (
        RouterLSA               uint8 = 1
        NetworkLSA              uint8 = 2
        Summary3LSA             uint8 = 3
        Summary4LSA             uint8 = 4
        ASExternalLSA           uint8 = 5
)

type LsaKey struct {
        LSType          uint8 /* LS Type */
        LSId            uint32 /* Link State Id */
        AdvRouter       uint32 /* Avertising Router */
}

func NewLsaKey() *LsaKey {
        return &LsaKey{}
}

type TOSDetail struct {
        TOS             uint8
        TOSMetric       uint16
}

type LinkDetail struct {
        LinkId          uint32 /* Link ID */
        LinkData        uint32 /* Link Data */
        LinkType        uint8 /* Link Type */
        NumOfTOS        uint8 /* # TOS Metrics */
        LinkMetric      uint16 /* Metric */
        TOSDetails      []TOSDetail
}

type LsaMetadata struct {
        LSAge           uint16 /* LS Age */
        Options         uint8 /* Options */
        LSSequenceNum   int /* LS Sequence Number */
        LSChecksum      uint16 /* LS Checksum */
        LSLen           uint16 /* LS Length */
}

/* LS Type 1 */
type RouterLsa struct {
        LsaMd           LsaMetadata
        BitV            bool /* V Bit */
        BitE            bool /* Bit E */
        BitB            bool /* Bit B */
        NumOfLinks      uint16 /* NumOfLinks */
        LinkDetails     []LinkDetail /* List of LinkDetails */
}

func NewRouterLsa() *RouterLsa {
        return &RouterLsa{}
}

/* LS Type 2 */
type NetworkLsa struct {
        LsaMd           LsaMetadata
        Netmask         uint32 /* Network Mask */
        AttachedRtr     []uint32
}

func NewNetworkLsa() *NetworkLsa {
        return &NetworkLsa{}
}

type SummaryTOSDetail struct {
        TOS             uint8
        TOSMetric       uint32
}

/* LS Type 3 or 4 */
type SummaryLsa struct {
        LsaMd                   LsaMetadata
        Netmask                 uint32 /* Network Mask */
        Metric                  uint32
        SummaryTOSDetails       []SummaryTOSDetail /* TOS */
}

func NewSummaryLsa() *SummaryLsa {
        return &SummaryLsa{}
}

type ASExtTOSDetail struct {
        BitE                    bool
        TOS                     uint8
        TOSMetric               uint32
        TOSFwdAddr              uint32
        TOSExtRouteTag          uint32

}

/* LS Type 5 */
type ASExternalLsa struct {
        LsaMd                   LsaMetadata
        Netmask                 uint32 /* Network Mask */
        BitE                    bool
        Metric                  uint32 /* But only max value 2^24-1 */
        FwdAddr                 uint32
        ExtRouteTag             uint32
        ASExtTOSDetails         []ASExtTOSDetail
}

func NewASExternalLsa() *ASExternalLsa {
        return &ASExternalLsa{}
}

type LSDatabase struct {
        RouterLsaMap            map[LsaKey]RouterLsa
        NetworkLsaMap           map[LsaKey]NetworkLsa
        Summary3LsaMap          map[LsaKey]SummaryLsa
        Summary4LsaMap          map[LsaKey]SummaryLsa
        ASExternalLsaMap        map[LsaKey]ASExternalLsa
}

type SelfOrigLsa map[LsaKey]bool
var InitialSequenceNumber int = 0x80000001
var MaxSequenceNumber int = 0x7fffffff
var LSSequenceNumber int = InitialSequenceNumber

/*

        0                   1                   2                   3
        0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |            LS age             |    Options    |    LS type    |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |                        Link State ID                          |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |                     Advertising Router                        |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |                     LS sequence number                        |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |         LS checksum           |             length            |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

*/

func encodeLsaHeader(lsaMd LsaMetadata, lsakey LsaKey) ([]byte) {
        lsaHdr := make([]byte, OSPF_LSA_HEADER_SIZE)
        binary.BigEndian.PutUint16(lsaHdr[0:2], lsaMd.LSAge)
        lsaHdr[2] = lsaMd.Options
        lsaHdr[3] = lsakey.LSType
        binary.BigEndian.PutUint32(lsaHdr[4:8], lsakey.LSId)
        binary.BigEndian.PutUint32(lsaHdr[8:12], lsakey.AdvRouter)
        binary.BigEndian.PutUint32(lsaHdr[12:16], uint32(lsaMd.LSSequenceNum))
        binary.BigEndian.PutUint16(lsaHdr[18:20], lsaMd.LSLen)
        return lsaHdr
}


/*
        0                   1                   2                   3
        0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |            LS age             |     Options   |       1       |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |                        Link State ID                          |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |                     Advertising Router                        |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |                     LS sequence number                        |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |         LS checksum           |             length            |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |    0    |V|E|B|        0      |            # links            |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |                          Link ID                              |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |                         Link Data                             |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |     Type      |     # TOS     |            metric             |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |                              ...                              |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |      TOS      |        0      |          TOS  metric          |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |                          Link ID                              |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |                         Link Data                             |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |                              ...                              |
*/

func decodeRouterLsa(data []byte, lsa *RouterLsa, lsakey *LsaKey) {
        lsa.LsaMd.LSAge = binary.BigEndian.Uint16(data[0:2])
        lsa.LsaMd.Options = uint8(data[2])
        lsakey.LSType = uint8(data[3])
        lsakey.LSId = binary.BigEndian.Uint32(data[4:8])
        lsakey.AdvRouter = binary.BigEndian.Uint32(data[8:12])
        lsa.LsaMd.LSSequenceNum = int(binary.BigEndian.Uint32(data[12:16]))
        lsa.LsaMd.LSChecksum = binary.BigEndian.Uint16(data[16:18])
        lsa.LsaMd.LSLen = binary.BigEndian.Uint16(data[18:20])
        if data[20] & 0x04 != 0 {
                lsa.BitV = true
        } else {
                lsa.BitV = false
        }
        if data[20] & 0x02 != 0 {
                lsa.BitE = true
        } else {
                lsa.BitE = false
        }
        if data[20] & 0x01 != 0 {
                lsa.BitB = true
        } else {
                lsa.BitB = false
        }
        lsa.NumOfLinks = binary.BigEndian.Uint16(data[22:24])
        lsa.LinkDetails = make([]LinkDetail, lsa.NumOfLinks)
        start := 24
        end := 0
        for i := 0; i < int(lsa.NumOfLinks); i++ {
                end = start+4
                lsa.LinkDetails[i].LinkId = binary.BigEndian.Uint32(data[start:end])
                start = end
                end = start+4
                lsa.LinkDetails[i].LinkData = binary.BigEndian.Uint32(data[start:end])
                start = end
                end = start+1
                lsa.LinkDetails[i].LinkType = uint8(data[start])
                start = end
                end = start+1
                lsa.LinkDetails[i].NumOfTOS = uint8(data[start])
                start = end
                end = start+4
                lsa.LinkDetails[i].LinkMetric = binary.BigEndian.Uint16(data[start:end])
                start = end
                lsa.LinkDetails[i].TOSDetails = make([]TOSDetail, lsa.LinkDetails[i].NumOfTOS)
                for j := 0; j < int(lsa.LinkDetails[i].NumOfTOS); j++ {
                        end = start+2
                        lsa.LinkDetails[i].TOSDetails[j].TOS = uint8(start)
                        start = end
                        end = start+2
                        lsa.LinkDetails[i].TOSDetails[j].TOSMetric = binary.BigEndian.Uint16(data[start:end])
                        start = end
                }
        }
}


func encodeLinkData(lDetail LinkDetail, length int) ([]byte) {
        lData := make([]byte, length)
        binary.BigEndian.PutUint32(lData[0:4], lDetail.LinkId)
        binary.BigEndian.PutUint32(lData[4:8], lDetail.LinkData)
        lData[8] = lDetail.LinkType
        lData[9] = lDetail.NumOfTOS
        binary.BigEndian.PutUint16(lData[10:12], lDetail.LinkMetric)
        start := 12
        end := 0
        for i := 0; i < int(lDetail.NumOfTOS); i++ {
                size := 4
                end = start + size
                lData[start] = lDetail.TOSDetails[i].TOS
                binary.BigEndian.PutUint16(lData[start+2:end], lDetail.TOSDetails[i].TOSMetric)
                start = end
        }
        return lData
}

func encodeRouterLsa(lsa RouterLsa, lsakey LsaKey)([]byte) {
        rtrLsa := make([]byte, lsa.LsaMd.LSLen)
        lsaHdr := encodeLsaHeader(lsa.LsaMd, lsakey)
        copy(rtrLsa[0:20], lsaHdr)
        var val uint8 = 0
        if lsa.BitV == true {
                val = val | 1 << 2
        }
        if lsa.BitE == true {
                val = val | 1 << 1
        }
        if lsa.BitB == true {
                val = val | 1
        }
        rtrLsa[20] = val
        binary.BigEndian.PutUint16(rtrLsa[22:24], lsa.NumOfLinks)

        start := 24
        end := 0
        for i := 0; i < int(lsa.NumOfLinks); i++ {
                size := 12 + 4 * lsa.LinkDetails[i].NumOfTOS
                end = start + int(size)
                linkData := encodeLinkData(lsa.LinkDetails[i], int(size))
                copy(rtrLsa[start:end], linkData)
                start = end
        }
        return rtrLsa
}

/*
        0                   1                   2                   3
        0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |            LS age             |      Options  |      2        |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |                        Link State ID                          |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |                     Advertising Router                        |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |                     LS sequence number                        |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |         LS checksum           |             length            |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |                         Network Mask                          |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |                        Attached Router                        |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |                              ...                              |


*/

func encodeNetworkLsa(lsa NetworkLsa, lsakey LsaKey) ([]byte) {
        nLsa := make([]byte, lsa.LsaMd.LSLen)
        lsaHdr := encodeLsaHeader(lsa.LsaMd, lsakey)
        copy(nLsa[0:20], lsaHdr)
        binary.BigEndian.PutUint32(nLsa[20:24], lsa.Netmask)
        numOfAttachedRtr := (int(lsa.LsaMd.LSLen) - OSPF_LSA_HEADER_SIZE - 4)/4
        start := 24
        for i := 0; i < numOfAttachedRtr; i++ {
                end := start + 4
                binary.BigEndian.PutUint32(nLsa[start:end], lsa.AttachedRtr[i])
                start = end
        }
        return nLsa
}

func decodeNetworkLsa(data []byte, lsa *NetworkLsa, lsakey *LsaKey) {
        lsa.LsaMd.LSAge = binary.BigEndian.Uint16(data[0:2])
        lsa.LsaMd.Options = uint8(data[2])
        lsakey.LSType = uint8(data[3])
        lsakey.LSId = binary.BigEndian.Uint32(data[4:8])
        lsakey.AdvRouter = binary.BigEndian.Uint32(data[8:12])
        lsa.LsaMd.LSSequenceNum = int(binary.BigEndian.Uint32(data[12:16]))
        lsa.LsaMd.LSChecksum = binary.BigEndian.Uint16(data[16:18])
        lsa.LsaMd.LSLen = binary.BigEndian.Uint16(data[18:20])
        lsa.Netmask = binary.BigEndian.Uint32(data[20:24])
        numOfAttachedRtr := (int(lsa.LsaMd.LSLen) - OSPF_LSA_HEADER_SIZE - 4)/4
        lsa.AttachedRtr = make([]uint32, numOfAttachedRtr)
        start := 24
        for i := 0; i < numOfAttachedRtr; i++ {
                end := start + 4
                lsa.AttachedRtr[i] = binary.BigEndian.Uint32(data[start:end])
                start  = end
        }
}

/*
        0                   1                   2                   3
        0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |            LS age             |     Options   |    3 or 4     |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |                        Link State ID                          |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |                     Advertising Router                        |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |                     LS sequence number                        |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |         LS checksum           |             length            |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |                         Network Mask                          |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |      0        |                  metric                       |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |     TOS       |                TOS  metric                    |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |                              ...                              |

*/

func encodeSummaryLsa(lsa SummaryLsa, lsakey LsaKey) ([]byte) {
        sLsa := make([]byte, lsa.LsaMd.LSLen)
        lsaHdr := encodeLsaHeader(lsa.LsaMd, lsakey)
        copy(sLsa[0:20], lsaHdr)
        binary.BigEndian.PutUint32(sLsa[20:24], lsa.Netmask)
        binary.BigEndian.PutUint32(sLsa[24:28], lsa.Metric)
        numOfTOS := (int(lsa.LsaMd.LSLen) - OSPF_LSA_HEADER_SIZE - 8)/8
        start := 28
        for i := 0; i < numOfTOS; i++ {
                end := start + 4
                var temp uint32
                temp = uint32(lsa.SummaryTOSDetails[i].TOS) << 24
                temp = temp | lsa.SummaryTOSDetails[i].TOSMetric
                binary.BigEndian.PutUint32(sLsa[start:end], temp)
                start = end
        }
        return sLsa
}

func decodeSummaryLsa(data []byte, lsa *SummaryLsa, lsakey *LsaKey) {
        lsa.LsaMd.LSAge = binary.BigEndian.Uint16(data[0:2])
        lsa.LsaMd.Options = uint8(data[2])
        lsakey.LSType = uint8(data[3])
        lsakey.LSId = binary.BigEndian.Uint32(data[4:8])
        lsakey.AdvRouter = binary.BigEndian.Uint32(data[8:12])
        lsa.LsaMd.LSSequenceNum = int(binary.BigEndian.Uint32(data[12:16]))
        lsa.LsaMd.LSChecksum = binary.BigEndian.Uint16(data[16:18])
        lsa.LsaMd.LSLen = binary.BigEndian.Uint16(data[18:20])
        lsa.Netmask = binary.BigEndian.Uint32(data[20:24])
        temp := binary.BigEndian.Uint32(data[24:28])
        lsa.Metric = 0x00ffffff | temp
        numOfTOS := (int(lsa.LsaMd.LSLen) - OSPF_LSA_HEADER_SIZE - 8)/8
        lsa.SummaryTOSDetails = make([]SummaryTOSDetail, numOfTOS)
        start := 28
        for i := 0; i < numOfTOS; i++ {
                end := start + 4
                temp = binary.BigEndian.Uint32(data[start:end])
                lsa.SummaryTOSDetails[i].TOS = uint8((0xff000000 | temp) >> 24)
                lsa.SummaryTOSDetails[i].TOSMetric = 0x00ffffff | temp
                start  = end
        }
}

/*
        0                   1                   2                   3
        0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |            LS age             |     Options   |      5        |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |                        Link State ID                          |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |                     Advertising Router                        |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |                     LS sequence number                        |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |         LS checksum           |             length            |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |                         Network Mask                          |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |E|     0       |                  metric                       |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |                      Forwarding address                       |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |                      External Route Tag                       |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |E|    TOS      |                TOS  metric                    |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |                      Forwarding address                       |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |                      External Route Tag                       |
       +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
       |                              ...                              |
*/

func encodeASExternalLsa(lsa ASExternalLsa, lsakey LsaKey) ([]byte) {
        eLsa := make([]byte, lsa.LsaMd.LSLen)
        lsaHdr := encodeLsaHeader(lsa.LsaMd, lsakey)
        copy(eLsa[0:20], lsaHdr)
        binary.BigEndian.PutUint32(eLsa[20:24], lsa.Netmask)
        var temp uint32
        if lsa.BitE == true {
                temp = temp | 0x80000000
        }
        temp = temp | lsa.Metric
        binary.BigEndian.PutUint32(eLsa[24:28], temp)
        binary.BigEndian.PutUint32(eLsa[28:32], lsa.FwdAddr)
        binary.BigEndian.PutUint32(eLsa[32:36], lsa.ExtRouteTag)
        numOfTOS := (int(lsa.LsaMd.LSLen) - OSPF_LSA_HEADER_SIZE - 16)/8
        start := 36
        for i := 0; i < numOfTOS; i++ {
                end := start + 4
                temp = 0
                if lsa.ASExtTOSDetails[i].BitE == true {
                        temp = temp | 0x80000000
                }
                temp = temp | uint32(lsa.ASExtTOSDetails[i].TOS) << 24 |
                        lsa.ASExtTOSDetails[i].TOSMetric
                binary.BigEndian.PutUint32(eLsa[start:end], temp)
                start = end
                end = start + 4
                binary.BigEndian.PutUint32(eLsa[start:end], lsa.ASExtTOSDetails[i].TOSFwdAddr)
                start = end
                end = start + 4
                binary.BigEndian.PutUint32(eLsa[start:end], lsa.ASExtTOSDetails[i].TOSExtRouteTag)
                start = end
        }
        return eLsa
}

func decodeASExternalLsa(data []byte, lsa *ASExternalLsa, lsakey *LsaKey) {
        lsa.LsaMd.LSAge = binary.BigEndian.Uint16(data[0:2])
        lsa.LsaMd.Options = uint8(data[2])
        lsakey.LSType = uint8(data[3])
        lsakey.LSId = binary.BigEndian.Uint32(data[4:8])
        lsakey.AdvRouter = binary.BigEndian.Uint32(data[8:12])
        lsa.LsaMd.LSSequenceNum = int(binary.BigEndian.Uint32(data[12:16]))
        lsa.LsaMd.LSChecksum = binary.BigEndian.Uint16(data[16:18])
        lsa.LsaMd.LSLen = binary.BigEndian.Uint16(data[18:20])
        lsa.Netmask = binary.BigEndian.Uint32(data[20:24])
        if data[24] == 0 {
                lsa.BitE = false
        } else {
                lsa.BitE = true
        }
        temp := binary.BigEndian.Uint32(data[24:28])
        lsa.Metric = 0x00ffffff | temp
        lsa.FwdAddr = binary.BigEndian.Uint32(data[28:32])
        lsa.ExtRouteTag = binary.BigEndian.Uint32(data[32:36])
        numOfTOS := (int(lsa.LsaMd.LSLen) - OSPF_LSA_HEADER_SIZE - 16)/8
        lsa.ASExtTOSDetails = make([]ASExtTOSDetail, numOfTOS)
        start := 36
        for i := 0; i < numOfTOS; i++ {
                end := start + 4
                temp = binary.BigEndian.Uint32(data[start:end])
                if temp & 0x80000000 != 0 {
                        lsa.ASExtTOSDetails[i].BitE = true
                } else {
                        lsa.ASExtTOSDetails[i].BitE = false
                }
                lsa.ASExtTOSDetails[i].TOS = uint8((temp & 0x7f000000) >> 24)
                lsa.ASExtTOSDetails[i].TOSMetric = 0x00ffffff | temp
                start = end
                end = start + 4
                lsa.ASExtTOSDetails[i].TOSFwdAddr = binary.BigEndian.Uint32(data[start:end])
                start = end
                end = start + 4
                lsa.ASExtTOSDetails[i].TOSExtRouteTag = binary.BigEndian.Uint32(data[start:end])
                start  = end
        }
}