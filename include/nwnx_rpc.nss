const string NWNX_PREFIX = "NWNX.";
const string RPC_PLUGIN_ID = "RPC";

const string RPC_PLUGIN_SEPARATOR = "!";

const string RPC_RESET_GENERIC = "RPC_RESET_GENERIC_";
const string RPC_BUILD_GENERIC = "RPC_BUILD_GENERIC_";
const string RPC_BUILD_GENERIC_STREAM = "RPC_BUILD_GENERIC_STREAM_"
const string RPC_PULL_GENERIC_STREAM = "RPC_PULL_GENERIC_STREAM_"

const string RPC_GET_INT = "RPC_GET_INT_";
const string RPC_SET_INT = "RPC_SET_INT_";
const string RPC_GET_BOOL = "RPC_GET_BOOL_";
const string RPC_SET_BOOL = "RPC_SET_BOOL_";
const string RPC_GET_FLOAT = "RPC_GET_FLOAT_";
const string RPC_SET_FLOAT = "RPC_SET_FLOAT_";
const string RPC_GET_STRING = "RPC_GET_STRING_";
const string RPC_SET_STRING = "RPC_SET_STRING_";
const string RPC_GET_GFF = "RPC_GET_GFF_";
const string RPC_SET_GFF = "RPC_SET_GFF_";

const int RPC_PARAM_2_DEFAULT = 0;
const int RPC_START_BUILD_GENERIC = 1;
const int RPC_END_BUILD_GENERIC = 2;

// CallAction
void RPCResetBuildGenericEx();
void RPCBuildGenericEx(string sClient, string sAction);
void RPCBuildGenericStreamEx(string sClient, string sAction);
int RPCPullGenericStreamEx();
int RPCGetIntEx(string sParam1, int nParam2 = RPC_PARAM_2_DEFAULT);
void RPCSetIntEx(string sParam1, int nValue, int nParam2 = RPC_PARAM_2_DEFAULT);
int RPCGetBoolEx(string sParam1, int nParam2 = RPC_PARAM_2_DEFAULT);
void RPCSetBoolEx(string sParam1, int bValue, int nParam2 = RPC_PARAM_2_DEFAULT);
float RPCGetFloatEx(string sParam1, int nParam2 = RPC_PARAM_2_DEFAULT);
void RPCSetFloatEx(string sParam1, float fValue, int nParam2 = RPC_PARAM_2_DEFAULT);
string RPCGetStringEx(string sParam1, int nParam2 = RPC_PARAM_2_DEFAULT);
void RPCSetStringEx(string sParam1, string sValue, int nParam2 = RPC_PARAM_2_DEFAULT);
object RPCRetrieveCampaignObjectEx(string sVarName, location coLocation, int nParam2 = RPC_PARAM_2_DEFAULT);
int RPCStoreCampaignObjectEx(string sVarName, object oObject, int nParam2 = RPC_PARAM_2_DEFAULT);

void RPCResetBuildGenericEx() {
	NWNXSetString(RPC_PLUGIN_ID, RPC_RESET_GENERIC, "", RPC_PARAM_2_DEFAULT, "");
}

void RPCBuildGenericEx(string sClient, string sAction) {
	NWNXSetString(RPC_PLUGIN_ID, RPC_BUILD_GENERIC, sClient, RPC_PARAM_2_DEFAULT, sAction);
}

void RPCBuildGenericStreamEx(string sClient, string sAction) {
	NWNXSetString(RPC_PLUGIN_ID, RPC_BUILD_GENERIC_STREAM, sClient, RPC_PARAM_2_DEFAULT, sAction);
}

void RPCPullGenericStreamEx() {
	return NWNXGetInt(RPC_PLUGIN_ID, RPC_PULL_GENERIC_STREAM, "", RPC_PARAM_2_DEFAULT);
}

int RPCGetIntEx(string sParam1, int nParam2 = RPC_PARAM_2_DEFAULT) {
	return NWNXGetInt(RPC_PLUGIN_ID, RPC_GET_INT, sParam1, nParam2);
}

void RPCSetIntEx(string sParam1, int nValue, int nParam2 = RPC_PARAM_2_DEFAULT) {
	NWNXSetInt(RPC_PLUGIN_ID, RPC_SET_INT, sParam1, nParam2, nValue);
}

int RPCGetBoolEx(string sParam1, int nParam2 = RPC_PARAM_2_DEFAULT) {
	return NWNXGetInt(RPC_PLUGIN_ID, RPC_GET_BOOL, sParam1, nParam2);
}

void RPCSetBoolEx(string sParam1, int bValue, int nParam2 = RPC_PARAM_2_DEFAULT) {
	NWNXSetInt(RPC_PLUGIN_ID, RPC_SET_BOOL, sParam1, nParam2, bValue == 1);
}

float RPCGetFloatEx(string sParam1, int nParam2 = RPC_PARAM_2_DEFAULT) {
	return NWNXGetFloat(RPC_PLUGIN_ID, RPC_GET_FLOAT, sParam1, nParam2);
}

void RPCSetFloatEx(string sParam1, float fValue, int nParam2 = RPC_PARAM_2_DEFAULT) {
	NWNXSetFloat(RPC_PLUGIN_ID, RPC_SET_FLOAT, sParam1, nParam2, fValue);
}

string RPCGetStringEx(string sParam1, int nParam2 = RPC_PARAM_2_DEFAULT) {
	return NWNXGetString(RPC_PLUGIN_ID, RPC_GET_STRING, sParam1, nParam2);
}

void RPCSetStringEx(string sParam1, string sValue, int nParam2 = RPC_PARAM_2_DEFAULT) {
	NWNXSetString(RPC_PLUGIN_ID, RPC_SET_STRING, sParam1, nParam2, sValue);
}

object RPCRetrieveCampaignObjectEx(string sVarName, location coLocation, int nParam2 = RPC_PARAM_2_DEFAULT) {
	return RetrieveCampaignObject(NWNX_PREFIX + RPC_PLUGIN_ID, RPC_GET_GFF + RPC_PLUGIN_SEPARATOR + IntToString(nParam2) + RPC_PLUGIN_SEPARATOR + sVarName, coLocation);
}

int RPCStoreCampaignObjectEx(string sVarName, object oObject, int nParam2 = RPC_PARAM_2_DEFAULT) {
	return StoreCampaignObject(NWNX_PREFIX + RPC_PLUGIN_ID, RPC_SET_GFF + RPC_PLUGIN_SEPARATOR + IntToString(nParam2) + RPC_PLUGIN_SEPARATOR + sVarName, oObject);
}

// NWNX*
int RPCGetInt(string sClient, string sParam1, int nParam2 = RPC_PARAM_2_DEFAULT);
void RPCSetInt(string sClient, string sParam1, int nValue, int nParam2 = RPC_PARAM_2_DEFAULT);
float RPCGetFloat(string sClient, string sParam1, int nParam2 = RPC_PARAM_2_DEFAULT);
void RPCSetFloat(string sClient, string sParam1, float fValue, int nParam2 = RPC_PARAM_2_DEFAULT);
string RPCGetString(string sClient, string sParam1, int nParam2 = RPC_PARAM_2_DEFAULT);
void RPCSetString(string sClient, string sParam1, string sValue, int nParam2 = RPC_PARAM_2_DEFAULT);
object RPCRetrieveCampaignObject(string sClient, string sVarName, location coLocation);
int RPCStoreCampaignObject(string sClient, string sVarName, object oObject);

int RPCGetInt(string sClient, string sParam1, int nParam2 = RPC_PARAM_2_DEFAULT) {
	return NWNXGetInt(RPC_PLUGIN_ID, sClient, sParam1, nParam2);
}

void RPCSetInt(string sClient, string sParam1, int nValue, int nParam2 = RPC_PARAM_2_DEFAULT) {
	NWNXSetInt(RPC_PLUGIN_ID, sClient, sParam1, nParam2, nValue);
}

float RPCGetFloat(string sClient, string sParam1, int nParam2 = RPC_PARAM_2_DEFAULT) {
	return NWNXGetFloat(RPC_PLUGIN_ID, sClient, sParam1, nParam2);
}

void RPCSetFloat(string sClient, string sParam1, float fValue, int nParam2 = RPC_PARAM_2_DEFAULT) {
	NWNXSetFloat(RPC_PLUGIN_ID, sClient, sParam1, nParam2, fValue);
}

string RPCGetString(string sClient, string sParam1, int nParam2 = RPC_PARAM_2_DEFAULT) {
	return NWNXGetString(RPC_PLUGIN_ID, sClient, sParam1, nParam2);
}

void RPCSetString(string sClient, string sParam1, string sValue, int nParam2 = RPC_PARAM_2_DEFAULT) {
	NWNXSetString(RPC_PLUGIN_ID, sClient, sParam1, nParam2, sValue);
}

object RPCRetrieveCampaignObject(string sClient, string sVarName, location coLocation) {
	return RetrieveCampaignObject(NWNX_PREFIX + RPC_PLUGIN_ID, sClient + RPC_PLUGIN_SEPARATOR + sVarName, coLocation);
}

int RPCStoreCampaignObject(string sClient, string sVarName, object oObject) {
	return StoreCampaignObject(NWNX_PREFIX + RPC_PLUGIN_ID, sClient + RPC_PLUGIN_SEPARATOR + sVarName, oObject);
}