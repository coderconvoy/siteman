function fold(caller){
    var sib = caller.nextElementSibling
    console.log("sib ==", sib)

    if (sib.style.display !== "none") {
        sib.style.display = "none";
    }else {
        sib.style.display = "";
    }
}

foldns = { };

function showFile(fname,caller){
    var box = document.getElementById("filebox");
    foldns.fname = fname;
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



