module.exports = function(grunt) {

  var jsFiles = ['Gruntfile.js',
                       'js/lib/Model.js',
                       'js/lib/InstanceBuilder.js',
                       'js/lib/ServerModel.js',
                       'js/lib/ClientModel.js',
                       'js/models/**/*.js',
                       'test/**/*.js'];

  grunt.initConfig({
    concat: {
      options: {
        stripBanners: true
      },
      clientjs: {
        src: [
           'js/lib/Model.js',
           'js/lib/InstanceBuilder.js',
           'js/lib/ClientModel.js',
           'js/models/**/common.js',
           'js/models/**/client.js'
        ],
        dest: 'bin/assets/client.js'
      },
      serverjs: {
        src: [
           'js/lib/Model.js',
           'js/lib/InstanceBuilder.js',
           'js/lib/ServerModel.js',
           'js/models/**/common.js',
           'js/models/**/server.js'
        ],
        dest: 'bin/server.js'
      }
    },
    jshint: {
      files: jsFiles,
      options: {
        // options here to override JSHint defaults
        globals: {
          jQuery: true,
          console: true,
          module: true,
          document: true
        }
      }
    },
    watch: {
      js: {
        files: jsFiles,
        tasks: ['jshint', 'buildjs']
      },
      go: {
        files: ['server/**/*.go'],
        tasks: ['shell:stop','buildgo','shell:start']
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

  grunt.loadNpmTasks('grunt-contrib-jshint');
  grunt.loadNpmTasks('grunt-contrib-watch');
  grunt.loadNpmTasks('grunt-contrib-concat');
  grunt.loadNpmTasks('grunt-shell');


  grunt.registerTask('start',['shell:stop','shell:start']);
  grunt.registerTask('buildjs',['concat:clientjs','concat:serverjs']);
  grunt.registerTask('buildgo',['shell:gobuild']);
  grunt.registerTask('build',['buildgo','buildjs','start']);
  grunt.registerTask('default', ['build','watch']);

};
