foldns = { };


function setPath(p,treepos){
    foldns.fname = p;
    foldns.treepos = treepos;
    document.getElementById("loc-p").innerHTML = p;
    
}

function fold(caller,path){
    document.getElementById("foldiv").style.display = "";
    document.getElementById("filediv").style.display = "none";
    setPath(path,caller);

    var sib = caller.nextElementSibling;
    console.log("sib ==", sib)

    if (sib.style.display !== "none") {
        sib.style.display = "none";
    }else {
        sib.style.display = "";
    }
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

function save(){
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
    fullname = foldns.fname + "/" + filename
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
                foldns.treepos.nextElementSibling.appendChild(nleaf);

            }else {
                console.log("No treepos",foldns.treepos);
            }
            setPath(fullname);
            showFile(fullname);

        }
    });
}


