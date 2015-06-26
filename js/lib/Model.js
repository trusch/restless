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
    if(typeof this.onPreInit === 'function'){
      this.onPreInit();
    }
    if(typeof this._doInit === 'function'){
      this._doInit(uid); //impl by subclass
    }
    if(typeof this.onPostInit === 'function'){
      this.onPostInit();
    }
    return this;
  }

  remove(){
    var res = false;
    if(typeof this.onPreRemove === 'function'){
      this.onPreRemove();
    }
    if(typeof this._doRemove === 'function'){
      res = this._doRemove(); //impl by subclass
    }
    if(typeof this.onPostRemove === 'function'){
      this.onPostRemove();
    }
    return res;
  }

  commit(){
    var res = false;
    if(typeof this.onPreCommit === 'function'){
      this.onPreCommit();
    }
    if(typeof this._doCommit === 'function'){
      res = this._doCommit(); //impl by subclass
    }
    if(typeof this.onPostCommit === 'function'){
      this.onPostCommit();
    }
    return res;
  }
}
