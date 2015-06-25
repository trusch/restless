class Model {

  constructor(){
    this.__uid = '';
    this.__rev = 0;
    this.__data = {};
    this.__schema = '';
  }

  initFromData(data){
    this.__data = data;
    return this;
  }

  initFromUID(uid){
    if(typeof this._preInit === 'function'){
      this._preInit();
    }
    if(typeof this._doInit === 'function'){
      this._doInit(uid); //impl by subclass
    }
    if(typeof this._postInit === 'function'){
      this._postInit();
    }
    return this;
  }

  remove(){
    var res = false;
    if(typeof this._preRemove === 'function'){
      this._preRemove();
    }
    if(typeof this._doRemove === 'function'){
      res = this._doRemove(); //impl by subclass
    }
    if(typeof this._postRemove === 'function'){
      this._postRemove();
    }
    return res;
  }

  commit(){
    var res = false;
    if(typeof this._preCommit === 'function'){
      this._preCommit();
    }
    if(typeof this._doCommit === 'function'){
      res = this._doCommit(); //impl by subclass
    }
    if(typeof this._postCommit === 'function'){
      this._postCommit();
    }
    return res;
  }
}
