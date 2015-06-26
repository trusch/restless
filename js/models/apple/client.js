
app.Model('Apple', CommonApple, class {
  onPreInit(){
    console.log('preInit');
  }
  onPreRemove(){
    console.log('preRemove');
  }
  onPreCommit(){
    console.log('preCommit');
  }
  onPostInit(){
    console.log('postInit');
  }
  onPostRemove(){
    console.log('postRemove');
  }
  onPostCommit(){
    console.log('postCommit');
  }
});

