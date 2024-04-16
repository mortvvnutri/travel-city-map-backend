const request = require('supertest');

app = request('http://localhost:5012');
const pre_path = '/api/public'
const genRanHex = size => [...Array(size)].map(() => Math.floor(Math.random() * 16).toString(16)).join('');

const access_denied = "403: Access denied"

const valid_email = "e1@example.com"
const valid_pwd = "yikes_long_password"


function basic_tests(url_list) {
  for (const url of url_list) {
    describe(`POST ${url}`, function() {
      it('Responds with json', function(done) {
        app.post(pre_path+url)
          .set('Accept', 'application/json')
          .send({})
          .expect('Content-Type', /json/)
          .expect(200, done);
      });
      it('Denies other methods', function(done) {
        app.get(pre_path+url)
          .expect(access_denied)
          .expect(403, done);
      });
    });
  }
}

basic_tests([
  '/categories/list',
  '/route/build',
  '/user/register',
  '/user/login',
  '/user/me',
])


describe('[SPECIFIC] POST /category/list', function() {
    it('Responds with some categories', function(done) {
      app.post(pre_path+'/category/list')
        .set('Accept', 'application/json')
        .send({})
        .expect(function(res) {
          res.body.categories.length>0
        })
        .expect(200, done);
    });
    it('Each category is valid', function(done) {
      app.post(pre_path+'/category/list')
        .set('Accept', 'application/json')
        .send({})
        .expect(function(res) {
          res.body.categories.every(cat=>{
            return cat.id!==undefined && cat.name!==undefined && cat.created_at!==undefined
          })
        })
        .expect(200, done);
    });
});

describe('[SPECIFIC] POST /category/list', function() {
  it('Responds with some categories', function(done) {
    app.post(pre_path+'/category/list')
      .set('Accept', 'application/json')
      .send({})
      .expect(function(res) {
        res.body.categories.length>0
      })
      .expect(200, done);
  });
  it('Each category is in valid format', function(done) {
    app.post(pre_path+'/category/list')
      .set('Accept', 'application/json')
      .send({})
      .expect(function(res) {
        res.body.categories.every(cat=>{
          return cat.id!==undefined && cat.name!==undefined && cat.created_at!==undefined
        })
      })
      .expect(200, done);
  });
  it('Each category has a valid existing parent', function(done) {
    app.post(pre_path+'/category/list')
      .set('Accept', 'application/json')
      .send({})
      .expect(function(res) {
        const parents = res.body.categories.filter(e=>e.parent_id===undefined)
        res.body.categories.every(e=>parents.includes(e.parent_id))
      })
      .expect(200, done);
  });
});

describe('[SPECIFIC] POST /user/register', function() {
  const url = '/user/register'
  const fake_email = () => {
    return `unit_test_${genRanHex(8)}@gmail.com`
  }
  const fake_pwd = genRanHex(20)
  const fake_name = `Unit Test ${genRanHex(8)}`
  it('Deny registration with incorrect data-types', function(done) {
    app.post(pre_path+url)
      .set('Accept', 'application/json')
      .send({
        user_initiator:{
          display_name: 1,
          pwd:3.33,
          email:{}
        }
      })
      .expect('Content-Type', /json/)
      .expect(res=>{
        res.body.error!==undefined
      })
      .expect(200, done);
  });
  it('Deny registration with short passwords', function(done) {
    app.post(pre_path+url)
      .set('Accept', 'application/json')
      .send({
        user_initiator:{
          display_name: fake_name,
          pwd:genRanHex(4),
          email:fake_email()
        }
      })
      .expect('Content-Type', /json/)
      .expect(res=>{
        res.body.error!==undefined
      })
      .expect(200, done);
  });
  it('Deny registration with short emails', function(done) {
    app.post(pre_path+url)
      .set('Accept', 'application/json')
      .send({
        user_initiator:{
          display_name: fake_name,
          pwd:genRanHex(14),
          email:"a@a.a"
        }
      })
      .expect('Content-Type', /json/)
      .expect(res=>{
        res.body.error!==undefined
      })
      .expect(200, done);
  });
  it('Deny registration with short names', function(done) {
    app.post(pre_path+url)
      .set('Accept', 'application/json')
      .send({
        user_initiator:{
          display_name: "d",
          pwd:fake_pwd,
          email:fake_email()
        }
      })
      .expect('Content-Type', /json/)
      .expect(res=>{
        res.body.error!==undefined
      })
      .expect(200, done);
  });
  it('Register a valid user, return token and userdata', function(done) {
    app.post(pre_path+url)
      .set('Accept', 'application/json')
      .send({
        user_initiator:{
          display_name: fake_name,
          pwd:fake_pwd,
          email:fake_email()
        }
      })
      .expect('Content-Type', /json/)
      .expect(res=>{
        res.body.error==undefined;
        res.body.token.length>32
        res.body.user.id!==undefined
        res.body.user.created_at!==undefined
      })
      .expect(200, done);
  });
});

describe('[SPECIFIC] POST /user/login', function() {
  const url = '/user/login'
  const fake_email = () => {
    return `unit_test_${genRanHex(8)}@gmail.com`
  }
  const fake_pwd = genRanHex(20)
  const fake_name = `Unit Test ${genRanHex(8)}`
  it('Deny login with incorrect data-types', function(done) {
    app.post(pre_path+url)
      .set('Accept', 'application/json')
      .send({
        user_initiator:{
          pwd:3.33,
          email:{}
        }
      })
      .expect('Content-Type', /json/)
      .expect(res=>{
        res.body.error!==undefined
      })
      .expect(200, done);
  });
  it('Deny login with non-existent credentials', function(done) {
      app.post(pre_path+url)
        .set('Accept', 'application/json')
        .send({
          user_initiator:{
            display_name: fake_name,
            pwd:genRanHex(40),
            email:fake_email()
          }
        })
        .expect('Content-Type', /json/)
        .expect(res=>{
          res.body.error!==undefined
        })
        .expect(200, done);
  });
  it('Login into a valid account, return token and userdata', function(done) {
    app.post(pre_path+url)
      .set('Accept', 'application/json')
      .send({
        user_initiator:{
          pwd:valid_pwd,
          email:valid_email
        }
      })
      .expect('Content-Type', /json/)
      .expect(res=>{
        res.body.error==undefined;
        res.body.token.length>32
        res.body.user.id!==undefined
        res.body.user.created_at!==undefined
      })
      .expect(200, done);
  });
});

describe('[SPECIFIC] POST /route/build', function() {
  const url = '/route/build'
  it('Deny route building with incorrect data-types', function(done) {
    app.post(pre_path+url)
      .set('Accept', 'application/json')
      .send({
        pos_req:{
          my_lat:"asd",
          my_long:"asd",
          cats:["a"]
        }
      })
      .expect('Content-Type', /json/)
      .expect(res=>{
        res.body.error!==undefined
      })
      .expect(200, done);
  });
  it('At least one category must be present', function(done) {
      app.post(pre_path+url)
        .set('Accept', 'application/json')
        .send({
          pos_req:{
            my_lat:37.6,
            my_long:55.7,
            cats:[]
          }
        })
        .expect('Content-Type', /json/)
        .expect(res=>{
          res.body.error!==undefined
        })
        .expect(200, done);
  });
  it("Don't include categories without places", function(done) {
      app.post(pre_path+url)
        .set('Accept', 'application/json')
        .send({
          pos_req:{
            my_lat:37.6,
            my_long:55.7,
            cats:[1,2,3,4,5,6]
          }
        })
        .expect('Content-Type', /json/)
        .expect(res=>{
          res.body.error==undefined
          res.body.places.length===0
        })
        .expect(200, done);
  });
  it("Include some categories many times if the user asks for it", function(done) {
      app.post(pre_path+url)
        .set('Accept', 'application/json')
        .send({
          pos_req:{
            my_lat:37.6,
            my_long:55.7,
            cats:[25,25,25,26,26,27]
          }
        })
        .expect('Content-Type', /json/)
        .expect(res=>{
          res.body.error==undefined
          res.body.places.filter(e=>e.id===25).length===3
          res.body.places.filter(e=>e.id===26).length===2
          res.body.places.filter(e=>e.id===27).length===1
        })
        .expect(200, done);
  });
  it("Ensure the best possible routing", function(done) {
      app.post(pre_path+url)
        .set('Accept', 'application/json')
        .send({
          pos_req:{
            my_lat:55.7287,
            my_long:37.6633,
            cats:[25,26,29,25]
          }
        })
        .expect('Content-Type', /json/)
        .expect(res=>{
          res.body.error==undefined
          JSON.stringify(res.body.places.map(e=>e.id))==="[37, 24900, 3154, 14143]"
        })
        .expect(200, done);
  });

});