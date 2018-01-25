let flist = {
    type : "folder"
}

exports.New = function (name){
        res = Object.create(flist)
        res.name = name;
        return res
}


