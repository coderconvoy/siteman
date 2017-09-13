foldns = { };
editns = {} ;

//Use as Callback for ajax edit errors
editns.error = function(resp) {
    showError("Edit Error: " + resp.responseText);
}

//Use as callback for ajax fs edits
editns.success = function(data){
    //Coming: Go through responses from server and enact basic operations
    for ( p in data) {
        var pp = data[p].Params;
        switch (data[p].Op.toLowerCase()) {
            case "say":
                showError(pp,true);
                break;
            case "err":
                showError(pp);
                break;
            case "rm" :
                editns.rm(pp)
                break;
            case "mkdir" : 
                editns.mkdir(pp);
                break;
            case "new" :
                editns.newFile(pp);
                break;
            case "mv" : 
                editns.mv(pp);
                break;
            default :
                showError("Unknown operation: " + data[p].Op);
                break;
        }
    }


}

editns.mkdirs = function(names) {
    var nn = names.split(",");
    for (p in nn) {
        editns.mkdir(nn[p]);
    }
}

editns.mkdir = function(fname) {
    var pp = fname.split("/");
    if (pp[0] == "") pp = pp.splice(0,1); 

    curr = document.getElementById("treetop");
    for (p in pp ){
        fileChild(pp[p])
    }
    //TODO
    
           /** if (teepos) {
                var nleaf = document.createElement("li");
                nleaf.innerHTML = folname;
                nleaf.onclick = function(){
                    fold(this);
                }
                nleaf.className = "treefolder";
                teepos.nextElementSibling.appendChild(nleaf);
                
                nchids = document.createElement("ul");
                teepos.nextElementSibling.appendChild(nchids);


            }else {
                console.log("No treepos",foldns.treepos);
            }
            setPath(fullname);
            showFile(nleaf);
            */
}


editns.newFile = function(fname){
}

editns.rm = function(fname){
}

editns.mv = function(fname){
}


function showError(mess,isHappy) {
    var emes = document.createElement("p");
    var eloc = document.getElementById("messbar");

    emes.innerHTML = mess;
    emes.className = "emess";
    if (isHappy){
        emes.className = "gmess";
    }
    emes.onclick = function(){
        this.remove();
        if (!eloc.firstChild){
            eloc.classList.add("hidden");
        }
        
    }

    eloc.appendChild(emes);
    eloc.classList.remove("hidden");
}

function setPath(p,treepos){
    foldns.fname = p;
    foldns.treepos = treepos;
    document.getElementById("loc-p").innerHTML = p;
    
}

function fold(caller){
    fpath = getPath(caller);
    console.log("Caller path:" + fpath);
    document.getElementById("foldiv").style.display = "";
    document.getElementById("filediv").style.display = "none";

    var sib = caller.nextElementSibling;

    if (sib.style.display !== "none" && foldns.treepos === caller) {
        sib.style.display = "none";
    }else {
        sib.style.display = "";
    }
    setPath(fpath,caller);
    var upcheck = document.getElementById("fup-location");
    upcheck.value = fpath;
}

function isIMG(fname){
    ext = fname.split('.').pop().toLowerCase();

    console.log("isimg:" + ext);
    
    switch (ext){
        case 'png':
        case 'gif':
        case 'svg':
        case 'jpg':
            return true;
    }
    return false;
}


function showFile(caller){
    var fname = getPath(caller);
    console.log("showFile: " , fname);
    document.getElementById("foldiv").style.display = "none";
    document.getElementById("filediv").style.display = "";

    var pic = document.getElementById("fileimg");
    var box = document.getElementById("filebox");
    setPath(fname,caller);

    var fpath = "/usr/"+fname;
    
    if (isIMG(fpath)) {
        pic.classList.remove("hidden");
        box.classList.add("hidden");
        pic.src = fpath;
    }else {
        pic.classList.add("hidden");
        box.classList.remove("hidden");
        box.value = "--LOADING--";
        $.get("/usr/"+fname,function(res){
            box.value = res ;
        });
    }
    console.log("Loading-" + fname)
}

function foldStart(){
    console.log("Hello fold starter");
}

function saveFile(){
    var fbox = document.getElementById("filebox");
    $.ajax({
        url:"/save",
        type:"POST",
        data:{
            fname:foldns.fname,
            fcontents:fbox.value
        },
        success:editns.success
    });
}

function addFolder(caller){
    var folname = prompt("Add Folder to " + foldns.fname + ": <br> Name :","untitled.txt");
    if (folname == null){
        return;
    }
    var fullname = foldns.fname + "/" + folname;
    var teepos = foldns.treepos;
    $.ajax({
        url:"/mkdir",
        type:"POST",
        data:{
            fname:fullname,
        },
        success:editns.success
    });
}

function addFile(caller){ 
     
    var filename = prompt("Add Folder to " + foldns.fname + ": <br> File Name :","untitled.txt");
    if (filename == null){
        return;
    }
    var fullname = foldns.fname + "/" + filename

    $.ajax({
        url:"/save",
        type:"POST",
        data:{
            fname:fullname,
            fcontents:""
        },
        success:function(){
            if (foldns.treepos) {
                var nleaf = document.createElement("li");
                nleaf.innerHTML = filename;
                nleaf.onclick = function(){
                    showFile(this);
                }
                nleaf.className = "treefile";
                foldns.treepos.nextElementSibling.appendChild(nleaf);

            }else {
                console.log("No treepos",foldns.treepos);
            }

            setPath(fullname);
            showFile(nleaf);

        },
        error:editns.error
    });
}

function postDelete(fname){
    $.ajax({
        url:"/delete",
        type:"POST",
        data:{
            fname:fname,
        },
        success:function(){
            showError("Deleted: " + fname,true);
            console.log("Deleting : ", foldns.treepos);
            foldns.treepos.remove();
        },
        error :function(resp){
            showError("Could not delete " + fname + ":" + resp.responseText);
            console.log(resp);
        }

    });
}

function deleteFolder(caller){
    ans = window.prompt("Delete whole folder: Are you sure? Please confirm by typing the full path of the folder you wish to delete.'"+foldns.fname+"'");
    if ( ans == "" ){
        showError("Delete Canceled: "+ foldns.fname); 
        return
    }
    if (ans !== foldns.fname){
        showError("Typed Incorrectly, not deleted: " + foldns.fname);
    }
    postDelete(foldns.fname);
}

function deleteFile(caller){
    if (!  confirm("Are you sure you want to delete " + foldns.fname+ "?")){
        showError("Delete Canceled: "+ foldns.fname); 
        return
    }
    postDelete(foldns.fname);
}

function selectFile(){
    var els = document.getElementsByClassName("with_select");
    for (var el in els){
        if (els[el].classList) {
        els[el].classList.remove("hidden");
        }
    }
    foldns.selectfname = foldns.fname;
    foldns.selectpos = foldns.treepos;

    
}

function deselectFile(){
    foldns.selectfname = undefined;
    foldns.selectpos = undefined;

    var els = document.getElementsByClassName("with_select");
    for (el in els){
        if (els[el].classList) {
        els[el].classList.add("hidden");
        }
    }
}

function getPath(treeitem){
    var path = treeitem.innerHTML;
    while(true){
        paritem = filePar(treeitem); 
        if (! paritem) {
            return path;
        }
        if (paritem.innerHTML == "/"){
            return path;
        }
        treeitem = paritem;
        path = treeitem.innerHTML + "/" + path;
    }
}

function filePar(pfile) {
    var paritem = pfile.parentNode.previousElementSibling;
    if (!paritem) {
        return undefined;
    }
    if (paritem.nodeName !== "LI") {
        return undefined;
    }
    return paritem;
}

function fileChild(pfile,childname) {
    plist = pfile.nextElementSibling;
    var cn = plist.children;
    if (!cn) {
        return undefined;
    }
    if (cn.nodeName !== "UL") {
        return undefined;
    }
    for ( var i = 0; i < cn.length; i++) {
        if (childname == cn[i].innerHTML) {
            return cn[i];
        }
    }
}

function descend(path){
    var root = document.getElementById("treetop");
    var curr = root;
    ps = path.split("/");
    for (var i =0 ; i < ps.length; i++ ) {
        if (i ==0 && ps[i] === "" ) continue;
        curr = fileChild(curr,ps[i]);
        if (!curr) {
            return undefined;
        }
    }
    return curr;
    
}

function moveHere(caller){
    
    newpath = foldns.fname + "/" + foldns.selectfname.split('/').pop();
    console.log("Moving file to " + newpath);
    
    $.ajax({
        url:"/move",
        type:"POST",
        data:{
            fname:foldns.selectfname,
            tname:newpath
        },
        success:function(data){
            console.log("Moving: " ,data);
            if (foldns.treepos && foldns.selectpos){
                var chids = foldns.selectpos.nextElementSibling;
                var newloc = foldns.treepos.nextElementSibling;
                newloc.appendChild(foldns.selectpos)
                if (chids ) if (chids.nodeName == "UL") { // if if for order
                    newloc.appendChild(chids); 
                }
            }
            deselectFile();
        },
        error:function(data){
            showError("ERROR:" +data);
        }

    });
}

function rename(caller){
    if (!foldns.treepos) {
        return;
    }
    var fname = foldns.treepos.innerHTML;
    tname = prompt("What would you like to rename '"+ fname+ "'?");

    if (!tname ) {
        return;
    }

    var fpath = getPath(foldns.treepos);
    var pathonly = fpath.substring(0,fpath.lastIndexOf("/"));
    console.log(",Pathonly : " ,pathonly);


    $.ajax({
        url:"/move",
        type:"POST",
        data:{
            fname:fpath,
            tname:pathonly + "/" +tname,
        },
        success:function(data){
            console.log("Moving: " ,data);
            if (foldns.treepos){
                foldns.treepos.innerHTML = tname;
            }
            deselectFile();
        },
        error:function(data){
            showError("ERROR:" +data);
        }
    });

    
}
