function fold(){
    var sib = this.nextSibling;

    if sib.style.display == "none" {
        sib.style.display = "";
    }else {
        sib.style.display = "none";
    }
}


