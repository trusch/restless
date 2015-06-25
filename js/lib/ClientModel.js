/*global InstanceBuilder*/
/*global Model*/

class ClientModel extends Model {
  constructor(){
    super();
  }

  _doInit(uid){
    console.log('ClientModel: init');
  }

  _doRemove(){
    console.log('ClientModel: remove');
  }
  
  _doCommit(){
    console.log('ClientModel: commit');
  }
}


