foldns = { };


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
    console.log("sib ==", sib)

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
        success:function(){
        }
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
        success:function(){
            if (teepos) {
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

        }
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
        error:function(){
            showError("Could not create file : " + nleaf);
        }
    });
}

function deleteFile(caller){
    if (!  confirm("Are you sure you want to delete " + foldns.fname+ "?")){
        return
    }

    $.ajax({
        url:"/delete",
        type:"POST",
        data:{
            fname:foldns.fname,
        },
        success:function(){
            showError("Deleted: " + foldns.fname,true);
            console.log("Deleting : ", foldns.treepos);
            foldns.treepos.remove();
        },
        error :function(){
            showError("Could not delete" + foldns.fname);
        }

    });

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
    path = treeitem.innerHTML;
    while(true){
        paritem = treeitem.parentNode.previousElementSibling;
        if (! paritem) {
            return path;
        }
        if (paritem.nodeName !== "LI"){
            return path;
        }
        if (paritem.innerHTML == "/"){
            return path;
        }
        treeitem = paritem;
        path = treeitem.innerHTML + "/" + path;
        
    }
}

function descend(path){
    var root = document.getElementById("treetop");
    console.log("root",root)
    var curr = root;
    bigloop:
    while (true){
        var cn = curr.children;
        for ( var i = 0; i < cn.length; i++) {
            if (path == cn[i].innerHTML) {
                return cn[i];
            }
            if (path.startsWith(cn[i].innerHTML + "/")) {
                console.log("in:"+cn[i].innerHTML);
                console.log("pre:" + path);
                path = path.slice(cn[i].innerHTML.length + 1);
                console.log("post:"+path);
                if (i +1 >= cn.length){
                    return undefined;
                }
                if(cn[i+1].nodeName !== "UL") {
                    return undefined;
                }
                curr = cn[i + 1];
                continue bigloop;
            }
            
        }
        return undefined ;
    }
    
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
