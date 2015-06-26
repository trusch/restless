/*global InstanceBuilder*/
/*global Model*/

class ServerModel extends Model {
  constructor(){
    super();
  }
  
  _doInit(uid){
    this.__uid = uid;
    let raw = db.get(this.__name+":"+this.__uid);
    this.__data = JSON.parse(raw);
  }

  _doRemove(){
    return db.remove(this.__name+":"+this.__uid);
  }
  
  _doCommit(){
    let key = this.__name+":"+this.__uid;
    let value = JSON.stringify(this.__data);
    console.log("_doCommit: ",key,value);
    return db.put(key,value);
  }
}


