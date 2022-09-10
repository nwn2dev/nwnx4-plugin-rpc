void RPCSetInt(string sClient, string sParam1, int nValue, int nParam2 = 0);
void RPCSetBool(string sClient, string sParam1, int nValue, bool bParam2 = 1);
void RPCSetFloat(string sClient, string sParam1, float fValue, int nParam2 = -1);
void RPCSetString(string sClient, string sParam1, string sValue, int nParam2 = -1);
int RPCGetInt(string sClient, string sParam1, int nParam2 = 0);
bool RPCGetBool(string sClient, string sParam1, int nParam2 = 1);
float RPCGetFloat(string sClient, string sParam1, int nParam2 = -1);
string RPCGetString(string sClient, string sParam1, int nParam2 = -1);
int RPCStoreCampaignObject(string sClient, string sVarName, object oObject);
object RPCRetrieveCampaignObject(string client, string sVarName, object oObject);

/*
Set an RPC int
*/
void RPCSetInt(string sClient, string sParam1, int nValue, int nParam2 = 0) {
	NWNXSetInt("RPC", sClient, sParam1, nParam2, nValue);
}

/*
Set an RPC bool
*/
void RPCSetBool(string sClient, string sParam1, bool bValue, int nParam2 = 1) {
	NWNXSetInt("RPC", sClient, sParam1, nParam2, bValue ? 1 : 0);
}

/*
Set an RPC float
*/
void RPCSetFloat(string sClient, string sParam1, float fValue, int nParam2 = -1) {
	NWNXSetFloat("RPC", sClient, sParam1, nParam2, fValue);
}

/*
Set an RPC string
*/
void RPCSetString(string sClient, string sParam1, string sValue, int nParam2 = -1) {
	NWNXSetString("RPC", sClient, sParam1, nParam2, fValue);
}

/*
Get an RPC int response
*/
int RPCGetInt(string sClient, string sParam1, int nParam2 = 0) {
	return NWNXGetInt("RPC", sClient, sParam1, nParam2);
}

/*
Get an RPC bool response
*/
bool RPCGetBool(string sClient, string sParam1, int nParam2 = 1) {
	return !(NWNXGetInt("RPC", sClient, sParam1, nParam2) == 0);
}

/*
Get an RPC float response
*/
float RPCGetFloat(string sClient, string sParam1, int nParam2 = -1) {
	return NWNXGetFloat("RPC", sClient, sParam1, nParam2);
}

/*
Get an RPC string response
*/
string RPCGetString(string sClient, string sParam1, int nParam2 = -1) {
	return NWNXGetString("RPC", sClient, sParam1, nParam2);
}

/*
Set a campaign object
*/
int RPCStoreCampaignObject(string sClient, string sVarName, object oObject) {
    return StoreCampaignObject("RPC", sClient + "///" + sVarName, oObject);
}

/*
Get a campaign object
*/
object RPCRetrieveCampaignObject(string client, string sVarName, object oObject) {
    return RetrieveCampaignObject("RPC", sClient + "///" + sVarName, GetLocation(oObject), oObject)
}
