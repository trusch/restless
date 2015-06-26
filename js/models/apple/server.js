/**
 * Apple Model: A simple model example
 */
app.Model('Apple', CommonApple, class {
  
  //onPreInit(){}
  onPostInit(){
    this.set('example.dynamic.requestTime',(new Date()).toString());
  }

  //onPreRemove(){}
  onPostRemove(){
    console.debug(`deleted instance ${this.__uid} of model ${this.__name}`);
  }

  onPreCommit(){
    this.set('lastModified',Date.now());
  }
  onPostCommit(){
    console.debug(`commited instance ${this.__uid} of model ${this.__name}`);
  }
  
  
});

