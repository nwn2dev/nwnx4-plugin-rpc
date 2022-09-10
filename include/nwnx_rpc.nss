int RPCGetInt(string sClient, string sParam1);
void RPCSetInt(string sClient, string sParam1, int nValue);

bool RPCGetBool(string sClient, string sParam1);
void RPCSetBool(string sClient, string sParam1, bool bValue);

float RPCGetFloat(string sClient, string sParam1);
void RPCSetFloat(string sClient, string sParam1, float fValue);

string RPCGetString(string sClient, string sParam1);
void RPCSetString(string sClient, string sParam1, string sValue);

object RPCRetrieveCampaignObject(string client, string sVarName, object oObject);
int RPCStoreCampaignObject(string sClient, string sVarName, object oObject);

const int RPC_TYPE_INT = 0;
const int RPC_TYPE_BOOL = 1;
const int RPC_TYPE_FLOAT = 2;
const int RPC_TYPE_STRING = 3;

const string RPC_FILE_SEPARATOR = "///";

/*
Set an RPC int
*/
void RPCSetInt(string sClient, string sParam1, int nValue) {
	NWNXSetInt("RPC", sClient, sParam1, RPC_TYPE_INT, nValue);
}

/*
Set an RPC bool
*/
void RPCSetBool(string sClient, string sParam1, bool bValue) {
	NWNXSetInt("RPC", sClient, sParam1, RPC_TYPE_BOOL, bValue ? 1 : 0);
}

/*
Set an RPC float
*/
void RPCSetFloat(string sClient, string sParam1, float fValue) {
	NWNXSetFloat("RPC", sClient, sParam1, RPC_TYPE_FLOAT, fValue);
}

/*
Set an RPC string
*/
void RPCSetString(string sClient, string sParam1, string sValue) {
	NWNXSetString("RPC", sClient, sParam1, RPC_TYPE_STRING, fValue);
}

/*
Get an RPC int response
*/
int RPCGetInt(string sClient, string sParam1) {
	return NWNXGetInt("RPC", sClient, sParam1, RPC_TYPE_INT);
}

/*
Get an RPC bool response
*/
bool RPCGetBool(string sClient, string sParam1) {
	return !(NWNXGetInt("RPC", sClient, sParam1, RPC_TYPE_BOOL) == 0);
}

/*
Get an RPC float response
*/
float RPCGetFloat(string sClient, string sParam1) {
	return NWNXGetFloat("RPC", sClient, sParam1, RPC_TYPE_FLOAT);
}

/*
Get an RPC string response
*/
string RPCGetString(string sClient, string sParam1) {
	return NWNXGetString("RPC", sClient, sParam1, RPC_TYPE_STRING);
}

/*
Set a campaign object
*/
int RPCStoreCampaignObject(string sClient, string sVarName, object oObject) {
    return StoreCampaignObject("RPC", sClient + RPC_FILE_SEPARATOR + sVarName, oObject);
}

/*
Get a campaign object
*/
object RPCRetrieveCampaignObject(string client, string sVarName, object oObject) {
    return RetrieveCampaignObject("RPC", sClient + RPC_FILE_SEPARATOR + sVarName, GetLocation(oObject), oObject)
}
