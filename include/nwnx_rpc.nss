const string NWNX_PREFIX = "NWNX.";
const string RPC_PLUGIN_ID = "RPC";

const string RPC_PLUGIN_SEPARATOR = "!";

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
const string RPC_RESET_CALL_ACTION = "RPC_RESET_CALL_ACTION_";
const string RPC_CALL_ACTION = "RPC_CALL_ACTION_";

const int RPC_CALL_ACTION_PARAM_2_DEFAULT = -1;
const int RPC_START_CALL_ACTION = 1;
const int RPC_END_CALL_ACTION = 2;

// CallAction
void RPCResetCallAction();
void RPCCallAction(string sClient, string sAction);
int RPCGetIntEx(string sParam1, int nParam2 = RPC_CALL_ACTION_PARAM_2_DEFAULT);
void RPCSetIntEx(string sParam1, int nValue, int nParam2 = RPC_CALL_ACTION_PARAM_2_DEFAULT);
int RPCGetBoolEx(string sParam1, int nParam2 = RPC_CALL_ACTION_PARAM_2_DEFAULT);
void RPCSetBoolEx(string sParam1, int bValue, int nParam2 = RPC_CALL_ACTION_PARAM_2_DEFAULT);
float RPCGetFloatEx(string sParam1, int nParam2 = RPC_CALL_ACTION_PARAM_2_DEFAULT);
void RPCSetFloatEx(string sParam1, float fValue, int nParam2 = RPC_CALL_ACTION_PARAM_2_DEFAULT);
string RPCGetStringEx(string sParam1, int nParam2 = RPC_CALL_ACTION_PARAM_2_DEFAULT);
void RPCSetStringEx(string sParam1, string sValue, int nParam2 = RPC_CALL_ACTION_PARAM_2_DEFAULT);
object RPCRetrieveCampaignObjectEx(string sVarName, int nParam2 = RPC_CALL_ACTION_PARAM_2_DEFAULT);
int RPCStoreCampaignObjectEx(string sVarName, object oObject, int nParam2 = RPC_CALL_ACTION_PARAM_2_DEFAULT);

void RPCResetCallAction() {
	NWNXSetString(RPC_PLUGIN_ID, RPC_RESET_CALL_ACTION, "", RPC_CALL_ACTION_PARAM_2_DEFAULT, "");
}

void RPCCallAction(string sClient, string sAction) {
	NWNXSetString(RPC_PLUGIN_ID, RPC_CALL_ACTION, sClient, RPC_CALL_ACTION_PARAM_2_DEFAULT, sAction);
}

int RPCGetIntEx(string sParam1, int nParam2 = RPC_CALL_ACTION_PARAM_2_DEFAULT) {
	return NWNXGetInt(RPC_PLUGIN_ID, RPC_GET_INT, sParam1, nParam2);
}

void RPCSetIntEx(string sParam1, int nValue, int nParam2 = RPC_CALL_ACTION_PARAM_2_DEFAULT) {
	NWNXSetInt(RPC_PLUGIN_ID, RPC_SET_INT, sParam1, nParam2, nValue);
}

int RPCGetBoolEx(string sParam1, int nParam2 = RPC_CALL_ACTION_PARAM_2_DEFAULT) {
	return NWNXGetInt(RPC_PLUGIN_ID, RPC_GET_BOOL, sParam1, nParam2);
}

void RPCSetBoolEx(string sParam1, int bValue, int nParam2 = RPC_CALL_ACTION_PARAM_2_DEFAULT) {
	NWNXSetInt(RPC_PLUGIN_ID, RPC_SET_BOOL, sParam1, nParam2, bValue == 1);
}

float RPCGetFloatEx(string sParam1, int nParam2 = RPC_CALL_ACTION_PARAM_2_DEFAULT) {
	return NWNXGetFloat(RPC_PLUGIN_ID, RPC_GET_FLOAT, sParam1, nParam2);
}

void RPCSetFloatEx(string sParam1, float fValue, int nParam2 = RPC_CALL_ACTION_PARAM_2_DEFAULT) {
	NWNXSetFloat(RPC_PLUGIN_ID, RPC_SET_FLOAT, sParam1, nParam2, fValue);
}

string RPCGetStringEx(string sParam1, int nParam2 = RPC_CALL_ACTION_PARAM_2_DEFAULT) {
	return NWNXGetString(RPC_PLUGIN_ID, RPC_GET_STRING, sParam1, nParam2);
}

void RPCSetStringEx(string sParam1, string sValue, int nParam2 = RPC_CALL_ACTION_PARAM_2_DEFAULT) {
	NWNXSetString(RPC_PLUGIN_ID, RPC_SET_STRING, sParam1, nParam2, sValue);
}

object RPCRetrieveCampaignObjectEx(string sVarName, int nParam2 = RPC_CALL_ACTION_PARAM_2_DEFAULT) {
	return RetrieveCampaignObject(NWNX_PREFIX + RPC_PLUGIN_ID, RPC_GET_GFF + RPC_PLUGIN_SEPARATOR + IntToString(nParam2) + RPC_PLUGIN_SEPARATOR + sVarName, GetLocation(OBJECT_SELF));
}

int RPCStoreCampaignObjectEx(string sVarName, object oObject, int nParam2 = RPC_CALL_ACTION_PARAM_2_DEFAULT) {
	return StoreCampaignObject(NWNX_PREFIX + RPC_PLUGIN_ID, RPC_SET_GFF + RPC_PLUGIN_SEPARATOR + IntToString(nParam2) + RPC_PLUGIN_SEPARATOR + sVarName, oObject);
}

// NWNX*
int RPCGetInt(string sClient, string sParam1, int nParam2 = 0);
void RPCSetInt(string sClient, string sParam1, int nValue, int nParam2 = 0);
float RPCGetFloat(string sClient, string sParam1, int nParam2);
void RPCSetFloat(string sClient, string sParam1, float fValue, int nParam2 = 0);
string RPCGetString(string sClient, string sParam1, int nParam2 = 0);
void RPCSetString(string sClient, string sParam1, string sValue, int nParam2 = 0);
object RPCRetrieveCampaignObject(string sClient, string sVarName);
int RPCStoreCampaignObject(string sClient, string sVarName, object oObject);

int RPCGetInt(string sClient, string sParam1, int nParam2 = 0) {
	return NWNXGetInt(RPC_PLUGIN_ID, sClient, sParam1, nParam2);
}

void RPCSetInt(string sClient, string sParam1, int nValue, int nParam2 = 0) {
	NWNXSetInt(RPC_PLUGIN_ID, sClient, sParam1, nParam2, nValue);
}

float RPCGetFloat(string sClient, string sParam1, int nParam2) {
	return NWNXGetFloat(RPC_PLUGIN_ID, sClient, sParam1, nParam2);
}

void RPCSetFloat(string sClient, string sParam1, float fValue, int nParam2 = 0) {
	NWNXSetFloat(RPC_PLUGIN_ID, sClient, sParam1, nParam2, fValue);
}

string RPCGetString(string sClient, string sParam1, int nParam2 = 0) {
	return NWNXGetString(RPC_PLUGIN_ID, sClient, sParam1, nParam2);
}

void RPCSetString(string sClient, string sParam1, string sValue, int nParam2 = 0) {
	NWNXSetString(RPC_PLUGIN_ID, sClient, sParam1, nParam2, sValue);
}

object RPCRetrieveCampaignObject(string sClient, string sVarName) {
	return RetrieveCampaignObject(NWNX_PREFIX + RPC_PLUGIN_ID, sClient + RPC_PLUGIN_SEPARATOR + sVarName, GetLocation(OBJECT_SELF));
}

int RPCStoreCampaignObject(string sClient, string sVarName, object oObject) {
	return StoreCampaignObject(NWNX_PREFIX + RPC_PLUGIN_ID, sClient + RPC_PLUGIN_SEPARATOR + sVarName, oObject);
}