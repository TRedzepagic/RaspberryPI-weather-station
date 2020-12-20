<template>
    <div class="hello">
      <h1>{{ msg }}</h1>
      <form >
              <div class="form-group">
          <label for="temperature">ID</label>
          <input
            type="number"
            placeholder="Ex: 14"
            v-model="form.id"
            class="form-control"
          />
        </div>
        <div class="form-group">
          <label for="temperature">Temperature</label>
          <input
            type="number"
            placeholder="Ex: 39.2"
            v-model="form.temperature"
            class="form-control"
          />
        </div>
        <div class="form-group" >
          <label for="pressure">Pressure</label>
          <input
            type="number"
            placeholder="Ex: 976.3"
            v-model="form.pressure"
            class="form-control"
            
          />
        </div>
          <div class="form-group">
          <label for="altitude">Altitude</label>
          <input
            type="number"
            placeholder="Ex: 324.5"
            v-model="form.altitude"
            class="form-control"
          />

        </div>
                <button type="button" v-on:click="postReading" class="btn btn-dark" :disabled=isDisabledPost>
          Add New Reading (POST)
        </button>
              <button type="button" v-on:click="updateReading" class="btn btn-dark" :disabled=isDisabledPut>
          Update Existing Reading (PUT)
        </button>
          <button type="button" v-on:click="getAll" class="btn btn-dark">
          Get All Readings (GET)
        </button>
              <button type="button" v-on:click="toggleLED" class="btn btn-dark">
          Toggle LED (RPi)
        </button>
        </form>


        <form v-on:submit.prevent>
        <div class="form-group">
          <label for="id">ID</label>
          <input
            type="number"
            placeholder="Ex: 14, THIS WILL DELETE THE READING FROM THE DATABASE!!!"
            v-model="formDelete.idDelete"
            class="form-control"
            required
          />
        </div>
        <input type="submit" value="Delete Reading (DELETE)" v-on:click="deleteReading" class="btn btn-danger" :disabled=isDisabledDelete>

      </form>

      
      <table class="table mt-5">
        <thead>
          <tr>
            <th scope="col">ID</th>
            <th scope="col">Temperature</th>
            <th scope="col">Pressure</th>
            <th scope="col">Altitude</th>
            <th scope="col">Time</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="reading in readingsList" :key="reading.id">
            <td>{{ reading.id }}</td>
            <td>{{ reading.temperature }}</td>
            <td>{{ reading.pressure }}</td>
            <td>{{ reading.altitude }}</td>
            <td>{{ reading.time }}</td>
          </tr>
        </tbody>
      </table>
      </div>
</template>

<script>
import axios from 'axios';
 var Timeout;
export default {
  name: 'main',
  props: { msg: String },
  data(){return {
    formDelete:{
      idDelete: '',
    },
    form:{
      id: '',
      temperature: '',
      pressure: '',
      altitude: '',
    },
    readingsList: [],
    errors: [],
    }},
  methods:{
    getAll() {
      axios.get(`http://192.168.0.23:8090/getReadings`)
      .then(response => {
      // JSON responses are automatically parsed.
      this.readingsList = response.data;
      console.log(response);
    })
    .catch(e => {
      console.log(e);
      this.errors.push(e);
      alert("GET errored out! "+e);
    })
     if(Timeout)
      {
        clearTimeout(Timeout);
      }
        Timeout = setTimeout(this.getAll, 5000);
  },



  postReading(){
    if(this.form.temperature=='' || this.form.pressure==''||this.form.altitude=='')
    {
      alert("One or more fields are empty, please correct this!");
      
    }else{
     axios.post(
    "http://192.168.0.23:8090/postReading",
    {
    temperature: this.form.temperature,
    pressure: this.form.pressure,
    altitude: this.form.altitude,
  },
  {
    headers: {
      "Content-type": "application/json; charset=UTF-8",
    }
  }).then(response => {
      // JSON responses are automatically parsed.
      console.log(response);
      alert("POST Successful!");
    })
    .catch(e => {
      console.log(e);
      this.errors.push(e);
      alert("POST errored out! "+e);
    }).finally(this.getAll)
    }
      if(Timeout)
      {
        clearTimeout(Timeout);
      }
    this.form.temperature='';
    this.form.pressure='';
    this.form.altitude='';
 
    
  },

  updateReading(){
     if(this.form.id==''||this.form.temperature=='' || this.form.pressure==''||this.form.altitude=='')
    {
      alert("One or more fields are empty, please correct this!");
      
    }else{
    axios.put(
    "http://192.168.0.23:8090/updateReading/"+this.form.id,
    {
    id: this.form.id,
    temperature: this.form.temperature,
    pressure: this.form.pressure,
    altitude: this.form.altitude,
  },
  {
    headers: {
      "Content-type": "application/json; charset=UTF-8",
    }
  }).then(response => {
      // JSON responses are automatically parsed.
      console.log(response);
      if(response.status==204)
      {
        alert("No reading with that ID");
      }else if(response.status==200){
      alert("PUT Successful!");
      }
    })
    .catch(e => {
      console.log(e);
      this.errors.push(e);
      alert("PUT errored out! "+e);
    }).finally(this.getAll)
      if(Timeout)
      {
        clearTimeout(Timeout);
      }
    this.form.id='';
    this.form.temperature='';
    this.form.pressure='';
    this.form.altitude='';
    }
  },

  deleteReading(){
    if (this.formDelete.idDelete ==''){
      alert("Please enter an ID!") 
      
    }else{
  axios.delete(
    "http://192.168.0.23:8090/deleteReading/"+this.formDelete.idDelete).then(response => {
      // JSON responses are automatically parsed.
      console.log(response);
      if(response.status==204)
      {
        alert("No reading with that ID");
      }else if(response.status==200){
      alert("DELETE Successful!");
      }
    })
    .catch(e => {
      console.log(e);
      this.errors.push(e);
      alert("DELETE errored out! "+e);
    }).finally(this.getAll)
      if(Timeout)
      {
        clearTimeout(Timeout);
      }
      this.formDelete.idDelete='';
    }
    
  },
  toggleLED(){
    axios.get(
    "http://RPI-IP:RPI-PORT/toggleLed/").then(response => {
      // JSON responses are automatically parsed.
      console.log(response);
      alert("LED Toggled! LED state: "+response.data.state);
    })
    .catch(e => {
      console.log(e);
      this.errors.push(e);
      alert("LED toggling errored out! "+e);
    })
  }, 


  },
  
  created(){
    this.getAll();
  },

  computed:{
  isDisabledPost(){
      if (this.form.id!='' || this.form.temperature=='' || this.form.pressue=='' || this.form.altitude==''){
      return true;
      }
      return false;
  },
  isDisabledPut(){
      if (this.form.id=='' || this.form.temperature=='' || this.form.pressue=='' || this.form.altitude==''){
      return true;
      }
      return false;
  },
  isDisabledDelete(){
      if (this.formDelete.idDelete==''){
      return true;
      }
      return false;
  },
  }  


}

</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h3 {
  margin: 40px 0 0;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
a {
  color: #42b983;
}

.btn{
  margin:1%
}

</style>
