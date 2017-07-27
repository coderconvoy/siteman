function fold(caller){
    console.log(caller);
    var sib = caller.nextElementSibling
    console.log("sib ==", sib)

    if (sib.style.visibilty !== "hidden") {
        sib.style.visibility = "hidden";
    }else {
        sib.style.display = "";
    }
}

function showFile(fname,caller){
    var box = document.getElementById("filebox");
    $.get("/usr/"+fname,function(res){
        box.value = res ;
    });
    console.log("Loading-" + fname)
}



function foldStart(){
    showFile("Hello Yall");
    console.log("Hello fold starter");
}



