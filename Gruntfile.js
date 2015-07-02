module.exports = function(grunt) {

  var jsFiles = [
     'js/lib/Aggregation.js',
     'js/lib/Model.js',
     'js/lib/ServerModel.js',
     'js/lib/ServerApp.js',
     'js/models/**/common.js',
     'js/models/**/server.js',
     'js/pages/**/*.js'
  ];

  grunt.initConfig({
    concat: {
      options: {
        stripBanners: true
      },
      serverjs: {
        src: jsFiles,
        dest: 'bin/server.js'
      }
    },
    babel: {
      options: {
        sourceMaps: true,
        comments: false
      },
      serverjs: {
        src: 'bin/server.js',
        dest: 'bin/server.js'
      }
    },
    watch: {
      js: {
        files: jsFiles,
        tasks: ['buildjs','start']
      },
      go: {
        files: ['server/**/*.go'],
        tasks: ['buildgo','start']
      }
    },
    shell: {
      gobuild: {
        command: 'cd server && go build -o ../bin/restless'
      },
      start: {
        command: 'echo ./bin/restless | at now'
      },
      stop: {
        command: 'killall restless || true'
      }
    }
  });

  grunt.loadNpmTasks('grunt-contrib-watch');
  grunt.loadNpmTasks('grunt-contrib-concat');
  grunt.loadNpmTasks('grunt-shell');
  grunt.loadNpmTasks('grunt-babel');


  grunt.registerTask('start',['shell:stop','shell:start']);
  grunt.registerTask('buildjs',['concat:serverjs','babel:serverjs']);
  grunt.registerTask('buildgo',['shell:gobuild']);
  grunt.registerTask('build',['buildgo','buildjs']);
  grunt.registerTask('default', ['build','watch']);

};
