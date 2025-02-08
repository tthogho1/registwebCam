package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"

	"registWebCam/util"
	"registWebCam/webcam"

	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	offset_max int = 50
)

var (
	baseurl    string = "https://api.windy.com/webcams/api/v3/webcams?lang=en&limit=" + strconv.Itoa(offset_max) + "&offset=%s&regions=%s"
	parameters string = "&sortDirection=asc&include=categories,images,location,player"
	logfile    string = "c:\\temp\\logrus.log"

	requestUrl string = baseurl + parameters
)

var logger = logrus.New()
var WINDY_API_KEY string

func main() {

	err := godotenv.Load(".env") // .envファイルを読み込む
	if err != nil {
		panic(err)
	}
	WINDY_API_KEY = os.Getenv("WINDY_API_KEY")

	// ファイル出力
	file, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		logger.Out = file
	} else {
		logger.Info("Failed to Open log to file, using default stderr")
	}

	regionCodes := extractRegionCode()
	//regionCodes := [...]string{"US.MD", "US.MT", "US.ID", "UY.02", "EC.01", "PT.10", "PT.23", "AR.01", "UY.10", "UY.09", "UY.08", "UY.07", "UY.16", "NZ.G2", "NZ.F3", "NZ.G3", "NZ.E9", "NZ.F7", "NZ.F8", "NZ.F2", "NZ.F9", "UY.06", "UY.05", "UY.17", "UY.04", "UY.14", "US.WI", "US.MN", "US.NV", "US.CA", "US.OR", "US.WA", "UY.03", "US.MI", "US.CT", "RS.SE", "MX.19", "MX.07", "UY.18", "UY.11", "US.HI", "UZ.02", "UZ.08", "TR.81", "MX.06", "MX.26", "AF.14", "AF.40", "AF.30", "AF.41", "AF.13", "AF.39", "AF.33", "CH.GL", "CH.LU", "CH.BE", "CH.AG", "CH.BL", "CH.ZG", "CH.OW", "CH.NW", "CH.GR", "CH.AR", "CH.AI", "CH.VS", "CH.SG", "CH.TI", "CH.SZ", "CH.ZH", "MD.78", "MD.57", "MD.59", "MD.82", "CH.TG", "MD.76", "MD.71", "MD.84", "CH.UR", "CH.SH", "CH.JU", "CH.FR", "MD.64", "CH.BS", "CH.SO", "CH.GE", "CH.NE", "CH.VD", "ET.45", "ET.46", "ET.51", "ET.54", "ET.52", "ET.44", "RW.12", "MA.02", "SO.03", "SO.06", "SO.14", "SO.18", "SO.12", "SO.20", "BG.38", "BG.39", "BG.62", "BG.63", "BG.64", "BG.41", "BG.40", "BG.44", "BG.45", "BG.46", "BG.48", "BG.49", "BG.50", "BG.51", "BG.52", "BG.53", "BG.55", "BG.58", "BG.43", "BG.65", "RU.42", "UZ.06", "UZ.03", "KG.09", "KG.08", "KG.06", "KG.03", "KG.07", "KG.02", "IT.03", "NZ.E8", "IN.29", "IT.06", "JP.08", "JP.38", "AO.01", "AO.02", "AO.09", "AO.20", "AO.14", "JP.03", "JP.33", "JP.27", "JP.25", "JP.21", "JP.18", "SY.01", "SY.11", "SY.04", "SY.09", "SY.10", "SY.12", "SY.02", "SY.08", "SY.06", "JP.10", "CU.02", "CU.03", "CU.01", "CU.05", "CU.07", "CU.09", "CU.10", "CU.13", "CU.15", "BR.16", "TR.84", "TR.25", "TR.04", "VN.77", "VN.73", "VN.65", "VN.67", "VN.93", "VN.87", "PE.23", "VN.21", "VN.01", "VN.09", "VN.69", "PE.18", "VN.03", "VN.37", "VN.24", "PE.04", "VN.55", "VN.49", "VN.88", "VN.91", "VN.23", "VN.60", "VN.54", "VN.61", "VN.46", "VN.63", "VN.84", "VN.78", "VN.66", "VN.62", "PE.06", "VN.52", "VN.58", "VN.34", "VN.76", "VN.33", "PE.11", "RU.91", "VN.59", "VN.82", "VN.81", "VN.13", "VN.79", "VN.74", "VN.86", "VN.83", "VN.71", "VN.30", "VN.85", "VN.53", "VN.70", "VN.32", "VN.90", "VN.68", "VN.50", "VN.44", "VN.47", "VN.45", "VN.43", "VN.75", "PE.21", "TH.35", "TH.28", "TH.44", "TH.32", "TH.26", "TH.48", "TH.02", "TH.03", "TH.46", "TH.58", "TH.23", "TH.11", "TH.50", "TH.22", "TH.63", "TH.06", "TH.05", "TH.18", "TH.34", "TH.01", "TH.24", "TH.43", "TH.53", "TH.73", "TH.27", "TH.16", "TH.64", "TH.04", "TH.31", "TH.38", "TH.39", "TH.69", "TH.61", "TH.66", "TH.41", "TH.14", "TH.56", "TH.13", "TH.12", "TH.36", "TH.07", "TH.74", "TH.57", "TH.59", "TH.52", "TH.47", "TH.25", "TH.42", "TH.54", "TH.37", "TH.67", "TH.33", "TH.30", "TH.68", "TH.09", "TH.51", "TH.60", "TH.29", "TH.08", "TH.65", "TH.49", "TH.75", "TH.76", "TH.10", "TH.70", "TH.72", "RU.31", "RU.87", "PE.25", "PE.08", "CU.14", "CU.16", "TH.79", "CD.08", "PE.05", "CD.19", "CD.27", "IN.07", "IN.36", "IN.10", "IN.23", "RO.01", "IN.24", "PE.07", "PE.15", "IN.09", "IN.35", "IN.16", "PE.02", "PE.10", "IN.34", "IN.28", "IN.38", "UZ.07", "UZ.12", "UZ.14", "UZ.16", "PE.13", "PE.14", "PE.22", "IN.37", "GG.8989934", "PE.01", "VN.20", "IN.21", "CN.05", "PE.20", "CL.01", "CL.12", "CN.08", "GE.04", "GE.65", "GE.72", "GE.71", "GE.51", "GE.66", "GE.69", "GE.67", "GE.68", "IN.13", "IN.19", "IN.02", "IN.03", "IN.26", "IN.30", "IN.14", "IN.18", "LY.69", "LY.79", "LY.80", "LY.65", "LY.72", "LY.73", "LY.67", "LY.76", "LY.77", "LY.68", "CL.08", "RO.34", "NZ.E7", "NZ.G1", "RU.77", "ML.07", "ML.04", "RU.49", "BW.11", "TW.02", "CN.31", "NZ.F6", "JP.22", "DE.02", "GT.01", "GT.12", "GT.14", "GT.16", "GT.07", "GT.09", "GT.10", "GT.17", "GT.08", "GT.13", "GT.15", "GT.19", "GT.06", "GT.03", "GT.11", "KZ.07", "KZ.04", "KZ.09", "KZ.15", "KZ.01", "KZ.17", "KZ.14", "KZ.10", "KZ.16", "KZ.11", "KZ.12", "CO.25", "ME.08", "NI.01", "NI.10", "NI.08", "NI.07", "KH.13", "KH.05", "KH.22", "KH.08", "KH.26", "KH.21", "KH.19", "KH.24", "KH.02", "UZ.13", "HR.10", "HR.18", "HR.02", "HR.17", "HR.11", "TR.31", "TR.32", "TR.07", "TR.71", "TR.33", "TR.48", "TR.78", "TR.46", "TR.83", "TR.02", "TR.09", "TR.35", "TR.45", "TR.20", "TR.64", "TR.43", "PL.73", "PL.85", "TR.11", "TR.26", "TR.68", "TR.60", "TR.58", "TR.44", "TR.63", "TR.21", "TR.72", "TR.76", "TR.49", "TR.13", "TR.65", "TR.23", "TR.12", "TR.24", "TR.77", "TR.61", "TR.08", "TR.28", "TR.52", "TR.05", "TR.19", "TR.82", "TR.89", "TR.37", "TR.57", "TR.87", "TR.85", "TR.93", "TR.14", "TR.54", "TR.16", "TR.10", "TR.22", "TR.39", "TR.59", "TR.17", "TR.92", "TR.41", "TR.34", "TR.62", "TR.69", "TR.53", "US.SC", "US.VA", "US.NC", "HR.14", "HR.01", "HR.09", "PL.72", "PL.74", "PL.77", "PL.81", "PL.83", "RO.27", "RO.33", "RO.09", "RO.15", "RO.20", "RO.28", "RO.04", "RO.02", "RO.21", "KN.07", "US.LA", "US.NJ", "IR.43", "RO.07", "RO.23", "RO.38", "RO.40", "RO.18", "RO.19", "RO.39", "RO.03", "HR.06", "HR.16", "HR.07", "HR.21", "HR.20", "HR.05", "HR.04", "HR.12", "HR.08", "HR.19", "HR.13", "HR.15", "HR.03", "VE.06", "VE.17", "VE.11", "VE.01", "NA.29", "NA.30", "NA.31", "NA.21", "NA.32", "NA.37", "NA.39", "IR.29", "KR.11", "KR.12", "IR.10", "IR.23", "IR.34", "IR.03", "KR.14", "KR.13", "KR.06", "AU.02", "AU.03", "AU.04", "AU.05", "AU.08", "CL.07", "AU.07", "ME.17", "ME.01", "ME.16", "TR.15", "ME.02", "ME.05", "ME.06", "ME.07", "ME.09", "ME.11", "ME.12", "ME.19", "ME.20", "ME.21", "KR.05", "KR.17", "MX.16", "MX.22", "MX.11", "MX.14", "MX.08", "KR.22", "KR.19", "AU.01", "KR.03", "RO.16", "RO.30", "RO.11", "RO.08", "RO.22", "RO.36", "RO.12", "RO.26", "RO.17", "RO.29", "RO.35", "RO.42", "RO.43", "RO.41", "RO.14", "RO.37", "AU.06", "CL.16", "CL.15", "ID.33", "ID.07", "ID.30", "ID.13", "ID.12", "ID.21", "ID.34", "ID.31", "ID.38", "ID.22", "ID.41", "ID.01", "ID.03", "ID.05", "ID.37", "ID.24", "ID.32", "ID.26", "KR.20", "KR.15", "KR.21", "KR.10", "ID.18", "ID.28", "ID.29", "KR.16", "CL.11", "KR.18", "MX.10", "AR.17", "CL.03", "IR.13", "IR.09", "IR.38", "IR.08", "IR.25", "IR.28", "IR.16", "MX.28", "MX.30", "MA.08", "MA.09", "MA.10", "MA.06", "MA.07", "BR.23", "MX.12", "VE.25", "MX.25", "KZ.02", "MQ.MQ", "GF.GF", "BH.19", "BH.17", "BA.BRC", "BA.01", "BA.02", "CL.06", "MX.20", "RU.19", "RU.01", "TD.23", "TD.28", "IR.07", "MX.05", "MX.04", "CU.AR", "MX.02", "MX.03", "MX.01", "MX.23", "MX.31", "CF.03", "CF.08", "NZ.F1", "BF.01", "JP.04", "JP.14", "JP.19", "PG.20", "IR.22", "IR.32", "IR.35", "IR.37", "MN.14", "MN.10", "MN.16", "MN.20", "BW.09", "BW.03", "CL.05", "VE.23", "VE.20", "VE.14", "VE.18", "VE.08", "VE.12", "VE.09", "VE.16", "VE.19", "VE.15", "VE.13", "VE.07", "VE.04", "QA.14", "QA.01", "RU.50", "HT.12", "BF.03", "TH.77", "CL.17", "CL.14", "CM.11", "CM.14", "CM.05", "CM.07", "CM.09", "HK.KKC", "DE.12", "BJ.09", "AR.09", "CN.16", "CN.18", "MZ.03", "MZ.10", "MZ.05", "JO.02", "JO.20", "JO.18", "JO.16", "JO.21", "JO.09", "JO.19", "JO.23", "MY.12", "TH.81", "TH.62", "TH.80", "TH.15", "MY.01", "MY.14", "MY.04", "MY.05", "KM.03", "KM.01", "KM.02", "CL.04", "BR.26", "IQ.11", "IQ.13", "IQ.15", "BR.18", "BR.27", "MZ.09", "BR.22", "BR.17", "IT.08", "CL.10", "BR.06", "BH.16", "BH.15", "BR.20", "BR.30", "BR.02", "BR.28", "CL.02", "EG.27", "EG.26", "EG.16", "EG.02", "EG.22", "EG.13", "EG.06", "EG.15", "AR.07", "KZ.05", "TD.13", "TD.05", "TD.06", "PG.19", "PG.16", "PG.09", "PG.02", "PG.06", "PG.18", "PG.11", "PG.12", "PG.14", "PG.03", "BR.15", "TN.23", "LB.04", "LB.08", "LB.05", "LB.06", "LB.07", "ZW.10", "ZM.02", "ZM.07", "ZM.09", "MG.7670842", "MG.7670848", "MG.7670854", "TD.19", "JP.11", "NE.01", "NE.07", "CR.04", "CR.03", "CR.01", "CR.08", "CR.02", "CR.06", "SL.03", "LK.38", "IQ.01", "IQ.10", "IQ.18", "IQ.04", "IQ.06", "IQ.02", "IQ.09", "IQ.12", "IQ.14", "IQ.03", "IQ.17", "IQ.16", "IN.40", "BR.24", "BR.01", "BR.25", "KH.28", "TJ.RR", "JP.41", "IR.40", "BR.03", "LR.19", "LR.18", "LR.09", "LR.11", "LR.20", "LR.14", "BR.04", "GH.08", "BR.13", "GH.02", "GH.04", "BR.14", "ZW.01", "ZW.04", "ZW.05", "ZW.08", "ZW.07", "ZW.02", "TG.22", "KE.11", "KE.41", "BR.11", "BR.29", "JP.36", "BO.08", "BR.31", "RU.66", "CI.93", "IR.11", "MR.07", "MR.01", "MR.11", "MR.02", "MR.14", "MR.06", "GN.F", "GN.K", "MA.12", "DO.15", "DO.04", "DO.08", "DO.25", "JP.32", "DO.30", "DO.31", "DO.23", "DO.32", "DO.21", "DO.01", "DO.35", "DO.10", "DO.20", "DO.37", "ID.08", "KE.39", "KE.19", "KE.21", "KE.46", "MU.12", "MU.13", "MU.14", "MU.15", "MU.16", "MU.17", "MU.18", "MU.19", "MU.20", "ES.07", "ES.32", "ES.34", "ES.58", "ES.55", "ES.59", "ES.51", "ES.52", "ES.31", "ES.53", "ES.57", "ES.54", "ES.56", "KE.31", "KE.32", "ZA.06", "ZA.02", "ZA.10", "ZA.09", "KE.18", "KE.22", "KE.23", "KE.27", "KE.28", "KE.29", "KE.30", "KE.37", "KE.44", "KE.45", "ZA.07", "KE.47", "KE.49", "KE.54", "GH.06", "GH.05", "GH.09", "EE.14", "EE.11", "EE.02", "EE.07", "EE.03", "EE.01", "EE.08", "EE.18", "EE.19", "AT.06", "EE.20", "EE.12", "EE.04", "EE.13", "EE.21", "EE.05", "JP.29", "JP.26", "JP.01", "KP.12", "CF.07", "KP.09", "JP.30", "JP.13", "JP.06", "JP.35", "JP.46", "CI.94", "PK.04", "EG.05", "AR.05", "CI.76", "CI.78", "CI.96", "CI.97", "SV.01", "SV.13", "SV.05", "BR.05", "SV.10", "SV.04", "SV.02", "SV.06", "SV.14", "SV.09", "SV.08", "NI.02", "NI.11", "AM.07", "AM.06", "AM.09", "AM.04", "AM.10", "AM.08", "AM.02", "AM.05", "AM.03", "AM.11", "IN.11", "BM.08", "BM.01", "BM.09", "BM.02", "BM.07", "IR.15", "SA.06", "SA.16", "ZM.10", "SA.15", "SA.13", "SA.10", "SA.11", "SA.14", "SA.05", "SA.08", "SA.02", "SA.17", "NG.44", "NG.23", "NG.45", "NG.22", "NG.36", "NG.37", "NG.54", "NG.47", "NG.11", "NG.05", "NG.31", "NG.16", "NG.43", "NG.27", "NG.46", "DJ.05", "DJ.07", "EG.04", "EG.18", "EG.10", "EG.28", "GQ.06", "FO.SU", "FO.VG", "FO.NO", "FO.ST", "FO.OS", "GN.04", "VU.13", "VU.16", "VU.18", "VU.15", "AE.01", "AE.02", "AE.03", "AE.04", "AE.05", "AE.07", "TO.01", "SD.49", "SD.50", "SD.56", "SD.61", "SD.36", "SD.29", "SD.38", "SD.43", "TW.01", "RO.10", "PG.13", "PG.15", "PK.03", "OM.03", "OM.06", "OM.11", "OM.01", "OM.09", "OM.10", "OM.08", "JP.24", "JP.44", "JP.02", "JP.16", "FR.44", "FR.84", "FR.27", "FR.75", "FR.76", "FR.28", "JP.37", "JP.15", "JP.31", "JP.39", "JP.20", "JP.05", "JP.17", "JP.47", "JP.12", "ID.35", "ID.40", "AZ.64", "FM.03", "FM.02", "EG.08", "EG.09", "EG.03", "BD.84", "BD.82", "CV.01", "CV.08", "CV.11", "SG.04", "YE.16", "YE.02", "YE.25", "YE.05", "YE.04", "YE.23", "YE.20", "YE.18", "YE.10", "YE.15", "YE.22", "YE.11", "YE.19", "YE.03", "YE.14", "YE.27", "SA.20", "BN.01", "BN.02", "BN.04", "BD.83", "HT.11", "MY.16", "MY.11", "SK.02", "SK.07", "SK.06", "SK.04", "SK.08", "SK.01", "SK.05", "SK.03", "PY.23", "PY.24", "PY.16", "PY.04", "PY.01", "PY.06", "PY.11", "CA.07", "CA.14", "CA.03", "CA.02", "CA.09", "CA.11", "CA.01", "CA.05", "CA.13", "CA.12", "IT.15", "AM.01", "HT.09", "BD.87", "BD.86", "BD.85", "BD.81", "US.RI", "BO.02", "RU.28", "RU.64", "BB.01", "BB.04", "BB.05", "BB.06", "BB.07", "BB.08", "BB.09", "BB.10", "BO.01", "BB.03", "CO.37", "MS.03", "NA.40", "BO.04", "IT.13", "IT.02", "VC.01", "VC.04", "CV.22", "CV.07", "CV.20", "CV.15", "CV.16", "CV.25", "CV.14", "IT.04", "NO.14", "NO.12", "NO.21", "NO.08", "CN.09", "LU.EC", "LU.GR", "LU.RM", "LU.ME", "LU.CA", "LU.LU", "IT.07", "NO.09", "MX.24", "EG.11", "EG.12", "EG.01", "EG.21", "EG.19", "EG.14", "IT.11", "IT.16", "IT.18", "BR.07", "FR.32", "AO.15", "IT.05", "NZ.F5", "NZ.F4", "NZ.TAS", "CZ.52", "IT.20", "MY.03", "MY.17", "MY.13", "MY.06", "MY.02", "MY.07", "MY.09", "GR.ESYE31", "IT.12", "IT.09", "BO.07", "IT.19", "ID.36", "ID.39", "IT.17", "KW.02", "KW.08", "KW.05", "KW.09", "ER.06", "HN.11", "HN.02", "HN.15", "HN.17", "HN.05", "HN.06", "HN.08", "HN.12", "HN.14", "HN.16", "HN.18", "LK.37", "LK.29", "LK.30", "LK.33", "LK.34", "TL.DI", "TL.MT", "FJ.01", "FJ.03", "FJ.05", "MA.04", "MA.03", "MA.05", "NL.02", "NL.16", "NL.01", "NL.03", "NL.15", "NL.07", "NL.09", "NL.06", "NL.11", "NL.05", "NL.10", "NL.04", "RU.47", "DE.10", "AT.07", "AT.02", "SE.14", "SE.23", "SE.07", "SE.24", "SE.03", "SE.10", "IT.10", "BE.VLG", "DO.34", "IR.42", "IR.04", "IR.39", "IT.01", "BE.BRU", "SE.21", "SE.25", "SE.15", "SE.22", "SE.28", "SE.08", "SE.18", "SE.26", "SE.06", "SE.27", "SE.12", "SE.02", "SE.09", "ID.14", "ID.42", "BR.08", "VN.39", "CN.02", "CN.07", "RO.06", "ID.10", "CD.11", "PA.04", "PA.07", "PA.13", "PA.08", "PA.10", "SR.16", "SR.18", "BR.21", "LA.03", "LA.15", "LA.17", "LA.19", "LA.20", "LA.24", "GB.WLS", "GB.SCT", "GB.ENG", "MD.60", "MD.92", "MD.80", "MD.77", "MD.91", "MD.85", "MD.73", "MD.63", "MD.81", "BY.02", "BY.06", "BY.04", "BY.03", "BY.07", "BY.05", "MM.03", "MM.16", "MM.04", "MM.15", "MM.08", "MM.13", "MM.01", "MM.10", "MM.11", "MM.12", "MM.17", "US.VT", "US.MA", "US.NY", "CA.10", "DE.08", "DE.15", "DE.09", "DE.16", "DE.13", "DE.11", "DE.14", "DE.01", "DE.05", "DE.07", "DE.06", "DE.04", "US.ME", "ID.04", "MD.67", "MD.88", "MV.43", "MV.39", "MV.44", "MV.31", "MV.30", "MV.10346475", "MV.41", "MV.32", "MV.33", "MV.34", "MV.38", "LU.RD", "KH.29", "US.NH", "CA.08", "CA.04", "MG.7670856", "MG.7670855", "MG.7670913", "MG.7670908", "MG.7670906", "UA.08", "GY.12", "UA.24", "GY.13", "UA.19", "UA.27", "UA.13", "UA.02", "UA.21", "UA.07", "RU.62", "UA.14", "UA.05", "UA.26", "TZ.23", "RU.43", "RU.86", "RU.57", "RU.81", "RU.67", "RU.65", "RU.51", "RU.46", "RU.83", "RU.41", "RU.56", "GM.02", "GM.07", "GM.05", "LB.10", "LB.11", "UA.15", "TO.EU", "RO.25", "RO.32", "UA.06", "UA.25", "UA.22", "UA.03", "UA.17", "UA.16", "UA.11", "BD.H", "BW.12", "TJ.04", "SY.13", "TT.03", "IT.14", "TT.CTT", "TT.PED", "TT.DMN", "TT.05", "TT.SJL", "TT.TUP", "TT.SGE", "TT.CHA", "GU.MA", "AT.08", "JP.09", "AT.01", "FR.94", "MX.18", "AT.03", "RU.84", "RU.55", "RU.08", "RU.13", "CK.11695425", "RU.73", "RU.71", "ZA.11", "RU.16", "JP.23", "RU.76", "RU.88", "RU.25", "RU.69", "RU.10", "RU.09", "NR.14", "NR.07", "KG.01", "GL.04", "GL.06", "GL.11839537", "GL.07", "RU.61", "RU.21", "MD.89", "RU.37", "FR.24", "FR.11", "FR.52", "AT.05", "FR.93", "ZA.08", "ZA.05", "YE.26", "RU.52", "BE.WAL", "UA.23", "UA.09", "JP.43", "CN.30", "UA.01", "UA.18", "CN.22", "CN.19", "CN.10", "CN.28", "CN.25", "CN.01", "CN.04", "CN.23", "CN.32", "CN.33", "CN.11", "CN.29", "CN.26", "CN.21", "CN.24", "CN.12", "CN.03", "RO.05", "RO.13", "TH.40", "ZA.03", "KN.08", "LC.08", "LC.03", "LC.05", "LC.06", "SE.16", "SE.05", "GT.18", "MD.62", "EC.08", "BZ.05", "BZ.01", "BZ.03", "BZ.02", "EC.07", "LT.64", "LT.59", "LT.65", "IN.39", "CK.11695124", "TD.26"}
	logger.Println(regionCodes)
	maxThreads := 2

	var wg sync.WaitGroup
	sem := make(chan struct{}, maxThreads)

	for _, regionCode := range regionCodes {
		wg.Add(1)
		sem <- struct{}{}

		go func(regionCode string) {
			defer wg.Done()
			defer func() { <-sem }()
			logger.Println(regionCode)
			extractAndRegistWebCamToMongoDB(regionCode)
		}(regionCode)

	}

	wg.Wait()
	logger.Println("All regions finished")
}

var (
	mongouri string = "mongodb+srv://webcam:webcam@cluster0.pizmgb2.mongodb.net/?retryWrites=true&w=majority"
)

func extractRegionCode() []string {
	url := "https://api.windy.com/webcams/api/v3/regions?lang=en"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("x-windy-api-key", WINDY_API_KEY)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	jsonArray := []map[string]interface{}{}
	err = json.Unmarshal(body, &jsonArray)
	if err != nil {
		panic(err)
	}

	codes := []string{}
	for _, item := range jsonArray {
		if item["code"] != nil {
			codes = append(codes, item["code"].(string))
		}
	}

	return codes
}

func extractWebCamData(url string) []byte {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	// ヘッダーを追加
	req.Header.Add("Accept", "application/json")
	req.Header.Add("x-windy-api-key", WINDY_API_KEY)

	// リクエストを送信
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	data, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	if res.StatusCode != 200 {
		fmt.Println("Status Code is " + strconv.Itoa(res.StatusCode) + " Body is " + string(data))
		log.Info("Status Code is " + strconv.Itoa(res.StatusCode) + " Body is " + string(data))
		return nil
	}
	res.Body.Close()

	return data
}

func registWebCamToMongoDB(data []byte, coll *mongo.Collection, ctx context.Context) {

	var webCameraInfo webcam.WebCameraInfo
	if err := json.Unmarshal(data, &webCameraInfo); err != nil {
		panic(err)
	}

	for _, webCam := range webCameraInfo.Webcams {

		result := coll.FindOne(context.TODO(), bson.M{"webcam.webcamid": webCam.WebcamID})
		if result.Err() != mongo.ErrNoDocuments {
			continue
		}

		var webCamWithEmd webcam.WebcamWithEmbedding
		webCamWithEmd.Webcam = webCam

		imgUrl := webCam.Images.Daylight.Thumbnail
		var imageData []byte
		if imgUrl != "" {
			// 画像をダウンロード
			imageData = util.GetImage(imgUrl)
			util.GetImage(imgUrl)
		}
		filename := strconv.Itoa(webCam.WebcamID) + ".jpg"
		print(filename)

		// upload to S3
		util.UploadS3(imageData, filename)

		rpgClient, _ := util.CreateClient()
		webCamWithEmd.Embedding = util.GetEmbedding(rpgClient, imageData, filename)

		_, err := coll.InsertOne(ctx, webCamWithEmd)
		if err != nil {
			print(err)
		}
		log.Info(webCam.WebcamID)
		fmt.Println(webCam.WebcamID)
	}
}

func extractAndRegistWebCamToMongoDB(regionCode string) {

	//  - バックグラウンドで接続する。タイムアウトは10秒
	ctx := context.TODO()

	// Create a new client and connect to the server
	con, err := mongo.Connect(ctx, options.Client().ApplyURI(mongouri))
	defer con.Disconnect(ctx)
	if err != nil {
		panic(err)
	}

	WEBCAM_DB := os.Getenv("WEBCAM_DB")
	WEBCAM_COLLECTION := os.Getenv("WEBCAM_COLLECTION")
	coll := con.Database(WEBCAM_DB).Collection(WEBCAM_COLLECTION)

	increment := 0
	for {
		// offset 取得
		offset := increment * offset_max
		increment++

		// url作成
		url := fmt.Sprintf(requestUrl, strconv.Itoa(offset), regionCode)
		fmt.Println(url)

		data := extractWebCamData(url)
		if data == nil {
			break
		}

		var webCameraInfo webcam.WebCameraInfo
		if err := json.Unmarshal(data, &webCameraInfo); err != nil {
			panic(err)
		}
		result_len := len(webCameraInfo.Webcams)
		//fmt.Println("Id : " + strconv.Itoa(id) + " data length : " + strconv.Itoa(result_len))

		registWebCamToMongoDB(data, coll, ctx)

		if result_len < offset_max {
			//fmt.Println("exit Id : " + strconv.Itoa(id) + " data length : " + strconv.Itoa(result_len))
			fmt.Println(webCameraInfo.Webcams)
			break
		}
	}
}
