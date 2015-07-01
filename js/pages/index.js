restless.Page('^/index\.js$',()=>{
    let pageData = fs.readFile('bin/assets/index.html')
    return {
        body: pageData+'\n\ndynamic: '+Date.now()
    };
});