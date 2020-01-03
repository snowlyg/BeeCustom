package com.hwt.bgcd.webservice;

import javax.xml.namespace.QName;

import org.apache.axiom.om.OMAbstractFactory;
import org.apache.axiom.om.OMElement;
import org.apache.axiom.om.OMFactory;
import org.apache.axiom.om.OMNamespace;
import org.apache.axis2.AxisFault;
import org.apache.axis2.addressing.EndpointReference;
import org.apache.axis2.client.Options;
import org.apache.axis2.rpc.client.RPCServiceClient;
import org.apache.axis2.transport.http.HTTPConstants;

public class Test {

	public static void main(String[] args) {		
//		System.out.println(requestACKMsg("MT2101A"));
		System.out.println(requestOnYardInfo("MRKU0925880"));		
//		send_report(Constants.xmlString);		
	}

	
	/**
	 * 发送预配舱单报文
	 * @param xmlString : 舱单报文
	 */
	private static void send_report(String xmlString){
		RPCServiceClient serviceClient = init_RPCService("http://www.cusdectrans.com:8014/BGCDWebService/services/OutBoundsService?wsdl");
		Object[] opAddEntryArgs = new Object[]{xmlString};
		Class[] classes = new Class[]{String.class};
		String result = null;
		QName opAddEntry = new QName("http://webservice.bgcd.hwt.com","out_bounds");
		try {			
			serviceClient.addHeader(createHeaderOMElement());
			result = serviceClient.invokeBlocking(opAddEntry,opAddEntryArgs, classes)[0].toString();
			System.out.println("KKKKKKKKKK:" + result);
		} catch (Exception e) {
			System.out.println(e.getMessage());
			e.printStackTrace();
		}		
	}
	
	/**
	 * 获取回执、船期、特殊费用
	 * 阻塞模式调用
	 * 参数(MT2101A、SDATE、SP_FEE)
	 * MT2101A：舱单回执
	 * SDATE:船期信息
	 * SP_FEE：特殊费用
	 * @return
	 */
	private static String requestACKMsg(String messageType) {
		RPCServiceClient serviceClient = init_RPCService("http://www.cusdectrans.com:8014/BGCDWebService/services/InBoundsService?wsdl");
		String result = null;
		Object[] opAddEntryArgs = new Object[]{messageType};
		Class[] classes = new Class[]{String.class}; 
		QName opAddEntry = new QName("http://webservice.bgcd.hwt.com","in_bounds");
		try {
			serviceClient.addHeader(createHeaderOMElement());
			result = serviceClient.invokeBlocking(opAddEntry,opAddEntryArgs, classes)[0].toString();	
	 	
		} catch (AxisFault e) {
			System.out.println(e.getMessage());
		}
		return result;
	}
	
	
	/**
	 * 获取在场箱信息
	 * @return
	 */
	private static String requestOnYardInfo(String conta_id){
		RPCServiceClient serviceClient = init_RPCService("http://www.cusdectrans.com:8014/BGCDWebService/services/QueryContaStatus?wsdl");
		String result = null;
		Object[] opAddEntryArgs = new Object[]{conta_id,"CNYTN"};
		Class[] classes = new Class[]{String.class}; 
		QName opAddEntry = new QName("http://webservice.bgcd.hwt.com","onYard_conta");
		try {			
			serviceClient.addHeader(createHeaderOMElement());
			result = serviceClient.invokeBlocking(opAddEntry,opAddEntryArgs, classes)[0].toString();
		} catch (AxisFault e) {
			System.out.println(e.getMessage());
		}		
		return result;
	}

	
	
	public static OMElement createHeaderOMElement(){ 
	  OMFactory factory = OMAbstractFactory.getOMFactory(); 
	     OMNamespace SecurityElementNamespace = factory.createOMNamespace("http://webservice.bgcd.hwt.com","authentication"); 
	        OMElement authenticationOM = factory.createOMElement("Authentication", 
	                SecurityElementNamespace); 
	        OMElement usernameOM = factory.createOMElement("Username", 
	                SecurityElementNamespace); 
	        OMElement passwordOM = factory.createOMElement("Password", 
	                SecurityElementNamespace); 
	        usernameOM.setText("ZSBS"); //登陆易报关系统用户名
	        passwordOM.setText("ZSBS123");//登陆易报关系统密码 
	        authenticationOM.addChild(usernameOM); 
	        authenticationOM.addChild(passwordOM); 
	        return authenticationOM; 
	 }
	
	private static RPCServiceClient init_RPCService(String url){
		//使用RPC方式调用WebService        
        RPCServiceClient serviceClient = null;
		try {
			serviceClient = new RPCServiceClient();
		} catch (AxisFault e) {
			System.out.println(e.getMessage());
		}
        Options options = serviceClient.getOptions();
        options.setProperty(HTTPConstants.SO_TIMEOUT, Integer.parseInt("3000"));//设置超时(单位是毫秒)        
        
        //指定调用WebService的URL
        EndpointReference targetEPR = new EndpointReference(url);
        options.setTo(targetEPR);
        return serviceClient;
	}
}
