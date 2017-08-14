foldns = { };


function setPath(p,treepos){
    foldns.fname = p;
    foldns.treepos = treepos;
    document.getElementById("loc-p").innerHTML = p;
    
}

function fold(caller,path){
    document.getElementById("foldiv").style.display = "";
    document.getElementById("filediv").style.display = "none";

    var sib = caller.nextElementSibling;
    console.log("sib ==", sib)

    if (sib.style.display !== "none" && foldns.treepos === caller) {
        sib.style.display = "none";
    }else {
        sib.style.display = "";
    }
    setPath(path,caller);
}


function showFile(fname,caller){
    document.getElementById("foldiv").style.display = "none";
    document.getElementById("filediv").style.display = "";

    var box = document.getElementById("filebox");
    setPath(fname,caller);
    $.get("/usr/"+fname,function(res){
        box.value = res ;
    });
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

function addFile(caller){ 
    
    var filename = document.getElementById("foltext").value;
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
                    showFile(fullname,this);
                }
                nleaf.className = "treefile";
                foldns.treepos.nextElementSibling.appendChild(nleaf);

            }else {
                console.log("No treepos",foldns.treepos);
            }
            setPath(fullname);
            showFile(fullname);

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
            console.log("Deleting : ", foldns.treepos);
            foldns.treepos.remove();
        }

    });

}

function selectFile(){
    var els = document.getElementsByClassName("with_select");
    for (el in els){
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
    
    $.ajax({
        url:"/move",
        type:"POST",
        data:{
            fname:foldns.selectfname,
            tname:newpath
        },
        success:function(){
            if (foldns.treepos ){
                foldns.treepos.nextElementSibling.appendChild(foldns.selectpos)
            }
            deselectFile();
        }

    });


}

