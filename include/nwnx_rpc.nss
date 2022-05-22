void RPCSetInt(string client, string sParam1, int nValue, int nParam2 = -1);
void RPCSetFloat(string client, string sParam1, float fValue, int nParam2 = -1);
void RPCSetString(string client, string sParam1, string sValue, int nParam2 = -1);
int RPCGetInt(string client, string sParam1, int nParam2 = -1);
float RPCGetFloat(string client, string sParam1, int nParam2 = -1);
string RPCGetString(string client, string sParam1, int nParam2 = -1);

/*
Set an RPC int
*/
void RPCSetInt(string client, string sParam1, int nValue, int nParam2 = -1) {
	NWNXSetInt("RPC", client, sParam1, nParam2, nValue);
}

/*
Set an RPC float
*/
void RPCSetFloat(string client, string sParam1, float fValue, int nParam2 = -1) {
	NWNXSetFloat("RPC", client, sParam1, nParam2, fValue);
}

/*
Set an RPC string
*/
void RPCSetString(string client, string sParam1, string sValue, int nParam2 = -1) {
	NWNXSetString("RPC", client, sParam1, nParam2, fValue);
}

/*
Get an RPC int response
*/
int RPCGetInt(string client, string sParam1, int nParam2 = -1) {
	return NWNXGetInt("RPC", client, sParam1, nParam2);
}

/*
Get an RPC float response
*/
float RPCGetFloat(string client, string sParam1, int nParam2 = -1) {
	return NWNXGetFloat("RPC", client, sParam1, nParam2);
}

/*
Get an RPC string response
*/
string RPCGetString(string client, string sParam1, int nParam2 = -1) {
	return NWNXGetString("RPC", client, sParam1, nParam2);
}
