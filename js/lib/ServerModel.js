/*global InstanceBuilder*/
/*global Model*/

class ServerModel extends Model {
  constructor(){
    super();
  }
  
  _doGet(uid){
    this.__uid = uid;
    let raw = db.get(this.__name+":"+this.__uid);
    this.__data = JSON.parse(raw);
  }

  _doRemove(){
    return db.remove(this.__name+":"+this.__uid);
  }
  
  _doPut(){
    let key = this.__name+":"+this.__uid;
    let value = JSON.stringify(this.__data);
    return db.put(key,value);
  }
}


