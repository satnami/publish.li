// --------------------------------------------------------------------------------------------------------------------

var app = new Vue({
  el   : '#app',
  data : {
    id      : '',
    title   : '',
    author  : '',
    content : '',
  },
  methods: {
    onSave: function() {
      axios.post('/', {
        id      : app.id,
        title   : app.title,
        author  : app.author,
        content : app.content,
      })
        .then(function (resp) {
          console.log(resp);
          var data = resp.data
          if ( data.ok ) {
            app.id = data.id
          }
        })
        .catch(function (err) {
          console.log(err);
        });
    },
  },
})

// --------------------------------------------------------------------------------------------------------------------
