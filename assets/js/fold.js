//Depends on treeview.js

foldns = {
    changed:false
};
editns = {} ;

debug = false;

//Use as Callback for ajax edit errors
editns.error = function(resp) {
    showError("Edit Error: " + resp.responseText);
}

//Use as callback for ajax fs edits
editns.success = function(data){
    //Coming: Go through responses from server and enact basic operations
    for ( p in data) {
        var pp = data[p].Params;
        if (debug) showError(data[p].Op +":"+ pp,true);
        switch (data[p].Op.toLowerCase()) {
            case "unchange":
                foldns.changed = false;
                break;
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
                editns.mkdirs(pp);
                break;
            case "new" :
                editns.newFiles(pp);
                break;
            case "mv" : 
                editns.mv.apply(editns,pp.split(","));
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

    curr = treeview.root(); 
    for (p in pp ){
        if (p == 0  && pp[0] == "") continue;

        var n = treeview.child(curr,pp[p]);
        if (!n ){
            n = treeview.addChildFolder(curr,pp[p],fold)   
            if (!n){
                showError("UI: Could not add child to " + curr.innerHTML);
                return undefined;
            }
        }
        curr = n;
    }
    return curr;
}


editns.newFiles = function(fname){
    console.log("newFiles: fname == " , fname);
    var ffs = fname.split(",");
    var res = undefined;
    for (p in ffs) {
        var path = ffs[p].split("/");
        var base = path.splice(-1,1);
        var pnode = editns.mkdir(path.join("/"));
        if (pnode) {
            res = treeview.addChildFile(pnode,base,showFile);
        }
    }
    showFile(res,true);
    return res;
}

editns.rm = function(fname){
    node = treeview.descend(fname);
    if (!node) {
        showError("Could not remove :" + fname + " : File not found");
        return
    }
    treeview.remove(node);
}

editns.mv = function(fname,tname){
    var node = treeview.descend(fname);
    if (!node) showError("File not found for move : " + fname);
    var blocks = treeview.remove(node);
    console.log(blocks);
    node.innerHTML = tname.split("/").pop();
    
    newloc = tname.substring(0,tname.lastIndexOf("/"));
    if (debug) showError("newloc = " + newloc,true);
    newpar = treeview.descend(newloc);
    treeview.addChildOb(newpar,blocks);

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
    if (foldns.changed ){
        if (! confirm(foldns.fname +" has been changed. Leave anyway?")){
            return;
        }
    }
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


function showFile(caller,override){
    if (foldns.changed && (!override)){
        if (! confirm(foldns.fname + " has been changed. Leave it anyway?")){
            return;
        }
    }
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
            foldns.changed = false;
            box.onchange = function(){foldns.changed = true;};
        });
    }
    console.log("Loading-" + fname)
}

function postAction(url,data,callback){
    if (!callback){
        callback = editns.success;
    }
    $.ajax({
        url:url,
        type:"POST",
        data:data,
        success:callback,
        error:editns.error
    });
}


function saveFile(){
    var fbox = document.getElementById("filebox");
    postAction("/save",{fname:foldns.fname,fcontents:fbox.value});
}

function addFolder(caller){
    var folname = prompt("Add Folder to " + foldns.fname + ": <br> Name :","untitled.txt");
    if (folname == null){
        return;
    }
    var fullname = foldns.fname + "/" + folname;
    var teepos = foldns.treepos;
    postAction("/mkdir",{fname:fullname});
}

function addFile(caller){ 
     
    var filename = prompt("Add Folder to " + foldns.fname + ": <br> File Name :","untitled.txt");
    if (filename == null){
        return;
    }
    var fullname = foldns.fname + "/" + filename

    postAction("/newfile",{fname:fullname,fcontents:""});
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
    postAction("/delete",{fname:foldns.fname});
}

function deleteFile(caller){
    if (!  confirm("Are you sure you want to delete " + foldns.fname+ "?")){
        showError("Delete Canceled: "+ foldns.fname); 
        return
    }
    postAction("/delete",{fname:foldns.fname});
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
        paritem = treeview.par(treeitem); 
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




function moveHere(caller){
    
    newpath = foldns.fname + "/" + foldns.selectfname.split('/').pop();
    console.log("Moving file to " + newpath);
    
    postAction("/move",{fname:foldns.selectfname,tname:newpath});
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

    postAction("/move",{fname:fpath,tname:pathonly + "/" + tname});
    
}
