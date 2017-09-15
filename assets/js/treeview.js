treeview = {};

treeview.root = function(){
    return document.getElementById("treetop").children[0];
}

treeview.descend = function(path){
    
    var curr = treeview.root();
    ps = path.split("/");
    console.log(ps);
    for (var i =0 ; i < ps.length; i++ ) {
        if (ps[i] === "" ) continue;
        console.log(curr);
        curr = treeview.child(curr,ps[i]);
        if (!curr) {
            return undefined;
        }
    }
    return curr;
}

treeview.par = function(pfile) {
    var paritem = pfile.parentNode.previousElementSibling;
    if (!paritem) {
        return undefined;
    }
    if (paritem.nodeName !== "LI") {
        return undefined;
    }
    return paritem;
}

treeview.child = function(pfile,childname) {
    plist = pfile.nextElementSibling;
    if (!plist) {
        console.log("No next sibling:",pfile);
    }
    if (plist.nodeName !== "UL") {
        console.log("not UL", plist, plist.nodeName); 
        return undefined;
    }
    var cn = plist.children;
    if (!cn) {
        console.log("no children", plist);
        return undefined;
    }
    for ( var i = 0; i < cn.length; i++) {
        if (childname == cn[i].innerHTML) {
            return cn[i];
        }
    }
    console.log("Child not in list :", childname, cn);
}

treeview.addChildOb = function(pfile,ob){
    var ul = pfile.nextElementSibling;
    if (! ul) {
        console.log("no sibling: ");console.log(pfile);
        return false;
    }
    if (ul.nodeName !== "UL") {
        showError("Could add folder under, " + curr.innerHTML);
        return false;
    }
    if (! Array.isArray(ob)){
        ob = [ob];
    }
    for (p in ob) {
        console.log(ob[p]);
        ul.appendChild(ob[p]);
    }
    return true; 
}

treeview.addChildFile = function(pfile,cname,clickf){
    var nleaf = document.createElement("li");
    nleaf.innerHTML = cname;
    nleaf.onclick = function(){
        clickf(this);
    }
    nleaf.className = "treefile";
    if (treeview.addChildOb(pfile,nleaf)) {
        return nleaf;
    }
    return undefined;

}

treeview.addChildFolder = function(pfile,cname,clickf){
    var nleaf = document.createElement("li");
    nleaf.innerHTML = cname;
    nleaf.onclick = function(){
        clickf(this);
    }
    nleaf.className = "treefolder";

    var nconts = document.createElement("ul");

    if (treeview.addChildOb(pfile,[nleaf,nconts]) ){
        return nleaf;
    }
    return undefined;
}

treeview.remove = function(pnode) {
    
    var fhold = pnode.nextElementSibling;
    pnode.remove();
    if (fhold) if (fhold.nodeName === "UL") {
        fhold.remove();
        return [pnode,fhold];
    }
    return pnode;
    

}


