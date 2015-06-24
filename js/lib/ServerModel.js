/*global InstanceBuilder*/
/*global Model*/

var ServerModel = {
  _doInit: function(uid){
    this.__uid = uid;
    var raw = db.get(this.__name+":"+this.__uid);
    this.__data = JSON.parse(raw);
  },
  _doRemove: function(){
    return db.remove(this.__name+":"+this.__uid);
  },
  _doCommit: function(){
    var key = this.__name+":"+this.__uid;
    var value = JSON.stringify(this.__data);
    console.log("_doCommit: ",key,value);
    return db.put(key,value);
  }
};

ServerModel = InstanceBuilder.extend(Model,ServerModel);

