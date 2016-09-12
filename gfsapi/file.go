package gfsapi

import (
	"strings"

	"github.com/Centny/gfs/gfsdb"
	"github.com/Centny/gwf/log"
	"github.com/Centny/gwf/routing"
	"github.com/Centny/gwf/util"
)

//List/Search User File/Folder
//List/Search login user file or folder
//
//@url,normal http get request
//	~/usr/api/listFile		GET
//@arg,the normal query arguments
//	name	O	the search key for file name
//	type	O	the type in `file/folder` to show the file or folder
//	pid		O	the parent folder id
//	tags	O	the file/folder tags to filter
//	pn		O	the page number begin of 1, default is 1
//	ps		O	the page size, default is 20
/*
	//
	//list user file or folder
	~/usr/api/listFile
	//list user file
	~/usr/api/listFile?type=file
	//searhc file
	~/usr/api/listFile?type=file&name=xx
*/
//@ret,code/data return
//	bases	O	the file base info, see upload api for deatail
//	files	A	the user file info.
/*	the example
	{
	    "code": 0,
	    "data": {
	        "bases": {
	            "57bd539dc3666e997e75f288": {
	                "exec": "done",
	                "ext": ".mp4",
	                "filename": "xx.mp4",
	                "id": "57bd539dc3666e997e75f288",
	                "info": {
	                    "V_json": {
	                        "count": 1,
	                        "files": ["www/2016-08-24/u_57bd539dc3666e997e000002_js.mp4"]
	                    },
	                    "V_pc": {
	                        "text": "www/2016-08-24/u_57bd539dc3666e997e000002_pc.mp4"
	                    },
	                    "code": 0
	                },
	                "md5": "52757d83284ca0967bc0c9e2be342c13",
	                "name": "xx.mp4",
	                "pub": "HI2hmt==",
	                "sha": "226cf3e82860ea778ccae40a9e424be5700249e1",
	                "size": 431684,
	                "status": "N",
	                "time": 1.472025501957e+12,
	                "type": "application/octet-stream"
	            }
	        },
	        "files": [{
	            "desc": "desc",
	            "fid": "57bd539dc3666e997e75f288",
	            "id": "57bd539dc3666e997e75f289",
	            "name": "xx.mp4",
	            "oid": "123",
	            "owner": "USR",
	            "pid": "57bd539ac3666e997e75f287",
	            "status": "N",
	            "tags": ["x", "y", "z"],
	            "time": 1.472025501961e+12,
	            "type": "file"
	        }]
	    }
	}
*/
//@tag,file,info,list
//@author,cny,2016-08-24
//@case,File System
func ListFile(hs *routing.HTTPSession) routing.HResult {
	var name, typ string
	var pid, ext, tags []string
	var pn, ps = 1, 20
	var err = hs.ValidCheckVal(`
		name,O|S,L:0;
		type,O|S,O:file~folder;
		ext,O|S,L:0;
		pid,O|S,L:0;
		tags,O|S,L:0;
		pn,O|I,R:0;
		ps,O|I,R:0;
		`, &name, &typ, &ext, &pid, &tags, &pn, &ps)
	if err != nil {
		return hs.MsgResErr2(1, "arg-err", err)
	}
	var uid = hs.StrVal("uid")
	if len(pid) < 1 {
		pid = []string{""}
	}
	fs, total, err := gfsdb.ListFilePaged(uid, OWN_USR, name, typ, pid, ext, tags, []string{gfsdb.FS_N}, pn-1, ps, 1)
	if err != nil {
		err = util.Err("ListFile list find by oid(%v),owner(%v),name(%v),type(%v),pid(%v),tags(%v) fail with error(%v)",
			uid, OWN_USR, name, typ, pid, tags, err)
		log.E("%v", err)
		return hs.MsgResErr2(2, "srv-err", err)
	}
	var ids = []string{}
	for _, f := range fs {
		ids = append(ids, f.Fid)
	}
	bases, err := gfsdb.ListF_m(ids)
	if err != nil {
		err = util.Err("ListFile list base file ids(%v) fail with error(%v)", ids, err)
		log.E("%v", err)
		return hs.MsgResErr2(3, "srv-err", err)
	}
	return hs.MsgRes(util.Map{
		"bases": bases,
		"files": fs,
		"total": total,
	})
}

//UpdateFile update user file or folder
//Update user file or foild by id
//
//@url,normal http get request
//	~/usr/api/updateFile?fid=xx		GET
//@arg,the normal query arguments
//	fid		R	the file/folder id
//	name	O	the file/folder name
//	desc	O	the file/folder desc
//	tags	O	the file/folder tags
/*
	//update file/folder name
	~/usr/api/updateFile?fid=xx&name=aaa
*/
//@ret,code/data return
//	code	I	the common code.
/*	the example
	{
	    "code": 0,
	    "data": "OK"
	}
*/
//@tag,file,info,update
//@author,cny,2016-08-24
//@case,File System
func UpdateFile(hs *routing.HTTPSession) routing.HResult {
	var file = &gfsdb.File{}
	var err = hs.ValidCheckVal(`
		fid,R|S,L:0;
		pid,O|S,L:0;
		name,O|S,L:0;
		desc,O|S,L:0;
		tags,O|S,L:0;
		`, &file.Id, &file.Pid, &file.Name, &file.Desc, &file.Tags)
	if err != nil {
		return hs.MsgResErr2(1, "arg-err", err)
	}
	file.Oid = hs.StrVal("uid")
	file.Owner = OWN_USR
	err = gfsdb.UpdateFile(file)
	if err != nil {
		return hs.MsgResErr2(2, "srv-err", err)
	}
	return hs.MsgRes("OK")
}

//RemoveFile remove user file or folder
//Remove user file or foild by id
//
//@url,normal http get request
//	~/usr/api/removeFile?fid=xx		GET
//@arg,the normal query arguments
//	fid		R	the file/folder id splited by comma
/*
	//remove file/folder
	~/usr/api/removeFile?fid=x1,x2
*/
//@ret,code/data return
//	code	I	the common code.
/*	the example
	{
	    "code": 0,
	    "data": "OK"
	}
*/
//@tag,file,remove
//@author,cny,2016-09-09
//@case,File System
func RemoveFile(hs *routing.HTTPSession) routing.HResult {
	var fid = hs.CheckVal("fid")
	if len(fid) < 1 {
		return hs.MsgResE3(1, "arg-err", "fid argument not found")
	}
	var err = gfsdb.RemoveFile(strings.Split(fid, ",")...)
	if err != nil {
		log.E("RemoveFile remove file by id(%v) fail with error(%v)", fid, err)
		return hs.MsgResErr2(2, "srv-err", err)
	}
	return hs.MsgRes("OK")
}

//AddFolder adding fild folder
//adding folder by name tags and parent foilder id
//
//@url,normal http get request
//	~/usr/api/addFolder?name=xx		GET
//@arg,the normal query arguments
//	pid		O	the parent folder id
//	name	O	the file/folder name
//	desc	O	the file/folder desc
//	tags	O	the file/folder tags
/*
	//adding folder
	~/usr/api/addFolder?name=aaa
*/
//@ret,code/data return
//	code	I	the common code.
//	data	I 	the added count, if return zero is meaning the folder is exists
/*	the example
	{
		"code": 0,
		"data": {
			"added": 1,
			"folder": {
				"desc": "",
				"fid": "",
				"id": "57d21be5c3666e08a4fe90d6",
				"name": "xx",
				"oid": "123",
				"owner": "USR",
				"pid": "",
				"status": "",
				"tags": [],
				"time": 0,
				"type": "folder"
			}
		}
	}
*/
//@tag,file,info,update
//@author,cny,2016-08-24
//@case,File System
func AddFolder(hs *routing.HTTPSession) routing.HResult {
	var file = &gfsdb.File{}
	var err = hs.ValidCheckVal(`
		pid,O|S,L:0;
		name,R|S,L:0;
		desc,O|S,L:0;
		tags,O|S,L:0;
		`, &file.Pid, &file.Name, &file.Desc, &file.Tags)
	if err != nil {
		return hs.MsgResErr2(1, "arg-err", err)
	}
	if len(file.Pid) > 0 {
		parent, err := gfsdb.FindFile(file.Pid)
		if err != nil {
			log.E("AddFolder find file by id(%v) fail with error(%v)", file.Pid, err)
			return hs.MsgResErr2(404, "srv-err", err)
		}
		if parent.Oid != hs.StrVal("uid") || parent.Owner != OWN_USR {
			log.E("AddFolder find file by id(%v) fail with error(%v)", file.Pid, "not access")
			return hs.MsgResErr2(404, "srv-err", util.Err("the parent folder is not found or not yours"))
		}
	}
	file.Oid, file.Owner, file.Type = hs.StrVal("uid"), OWN_USR, gfsdb.FT_FOLDER
	file.Status, file.Time = gfsdb.FS_N, util.Now()
	updated, err := gfsdb.FOI_File(file)
	if err != nil {
		log.E("AddFolder find file by id(%v) fail with error(%v)", file.Pid, err)
		return hs.MsgResErr2(2, "srv-err", err)
	}
	return hs.MsgRes(util.Map{
		"added":  updated,
		"folder": file,
	})
}
