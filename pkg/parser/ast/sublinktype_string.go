// Code generated by "stringer -type SubLinkType nodes/sub_link_type.go"; DO NOT EDIT.

package pg_query

import "strconv"

const _SubLinkType_name = "EXISTS_SUBLINKALL_SUBLINKANY_SUBLINKROWCOMPARE_SUBLINKEXPR_SUBLINKMULTIEXPR_SUBLINKARRAY_SUBLINKCTE_SUBLINK"

var _SubLinkType_index = [...]uint8{0, 14, 25, 36, 54, 66, 83, 96, 107}

func (i SubLinkType) String() string {
	if i >= SubLinkType(len(_SubLinkType_index)-1) {
		return "SubLinkType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _SubLinkType_name[_SubLinkType_index[i]:_SubLinkType_index[i+1]]
}
