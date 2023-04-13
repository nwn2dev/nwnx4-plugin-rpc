const string RPC_PLUGIN_ID = "RPC";

const string RPC_GFF_VAR_NAME_SEPARATOR = "///";

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

// CallAction
void RPCResetCallAction();
void RPCCallAction(string sClient, string sAction);
int RPCGetIntEx(string sParam1);
void RPCSetIntEx(string sParam1, int nValue);
int RPCGetBoolEx(string sParam1);
void RPCSetBoolEx(string sParam1, int bValue);
float RPCGetFloatEx(string sParam1);
void RPCSetFloatEx(string sParam1, float fValue);
string RPCGetStringEx(string sParam1);
void RPCSetStringEx(string sParam1, string sValue);
object RPCRetrieveCampaignObjectEx(string sVarName);
int RPCStoreCampaignObjectEx(string sVarName, object oObject);

void RPCResetCallAction() {
	NWNXSetString(RPC_PLUGIN_ID, RPC_RESET_CALL_ACTION, "", -1, "");
}

void RPCCallAction(string sClient, string sAction) {
	NWNXSetString(RPC_PLUGIN_ID, RPC_CALL_ACTION, sClient, -1, sAction);
}

int RPCGetIntEx(string sParam1) {
	return NWNXGetInt(RPC_PLUGIN_ID, RPC_GET_INT, sParam1, -1);
}

void RPCSetIntEx(string sParam1, int nValue) {
	NWNXSetInt(RPC_PLUGIN_ID, RPC_SET_INT, sParam1, -1, nValue);
}

int RPCGetBoolEx(string sParam1) {
	return NWNXGetInt(RPC_PLUGIN_ID, RPC_GET_BOOL, sParam1, -1);
}

void RPCSetBoolEx(string sParam1, int bValue) {
	NWNXSetInt(RPC_PLUGIN_ID, RPC_SET_BOOL, sParam1, -1, bValue != 0);
}

float RPCGetFloatEx(string sParam1) {
	return NWNXGetFloat(RPC_PLUGIN_ID, RPC_GET_FLOAT, sParam1, -1);
}

void RPCSetFloatEx(string sParam1, float fValue) {
	NWNXSetFloat(RPC_PLUGIN_ID, RPC_SET_FLOAT, sParam1, -1, fValue);
}

string RPCGetStringEx(string sParam1) {
	return NWNXGetString(RPC_PLUGIN_ID, RPC_GET_STRING, sParam1, -1);
}

void RPCSetStringEx(string sParam1, string sValue) {
	NWNXSetString(RPC_PLUGIN_ID, RPC_SET_STRING, sParam1, -1, sValue);
}

object RPCRetrieveCampaignObjectEx(string sVarName) {
	return RetrieveCampaignObject(RPC_PLUGIN_ID, RPC_GET_GFF + RPC_GFF_VAR_NAME_SEPARATOR + sVarName, GetStartingLocation());
}


int RPCStoreCampaignObjectEx(string sVarName, object oObject) {
	return StoreCampaignObject(RPC_PLUGIN_ID, RPC_SET_GFF + RPC_GFF_VAR_NAME_SEPARATOR + sVarName, oObject);
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
	return RetrieveCampaignObject(RPC_PLUGIN_ID, sClient + RPC_GFF_VAR_NAME_SEPARATOR + sVarName, GetStartingLocation());
}

int RPCStoreCampaignObject(string sClient, string sVarName, object oObject) {
	return StoreCampaignObject(RPC_PLUGIN_ID, sClient + RPC_GFF_VAR_NAME_SEPARATOR + sVarName, oObject);
}
