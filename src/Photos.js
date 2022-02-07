import React, { useState, useEffect} from 'react';
import { Button,Grid,FormControl,InputLabel,Select,MenuItem} from '@material-ui/core';
import { ArrowUpward } from '@material-ui/icons';
import MinorHeader from "./components/MinorHeader.js"
import './App.css';
import GuerNav from "./components/GuerNav.js"
import Moment from 'react-moment';
import 'moment-timezone';
import { Cancel } from '@material-ui/icons';


function Photos(props) {
const [filename, setFilename] = useState('');
const [filehash, setFilehash] = useState('');
const [photocount, setPhotocount] = useState('');
const [photos, setPhotos] = useState([]);
const [photoloading, setPhotoloading] = useState(false);
const [checkedphotos, setCheckedphotos] = useState(false);
const [filetype, setFiletype] = React.useState('');
let temparray = [];

useEffect(() => {
    if (!photocount && props.guer && props.myaccount && props.mymia){
        const photoload = async () => {
            const allphotos = await props.guer.methods.doGetPhotoCount(props.mymia).call({from: props.myaccount});
            setPhotocount(allphotos);
        }
        photoload();
    }
    if (checkedphotos === false && photocount > 0 && (temparray.length === 0 || temparray.length === "0")) {
        setCheckedphotos(true);
        const photodetails = async () => {
        for (let i = 1;i <= photocount;i++){
            await props.guer.methods.doGetPhoto(props.mymia, i).call({from: props.myaccount})
            .then(function(result){
                temparray.push(result);
            });
        }
        setPhotos(temparray);
        }
        photodetails();
    }
},[props.guer,props.myaccount,props.mymia,photocount,temparray,checkedphotos]);

const uploadFile = async (e) => {
    setPhotoloading(true);
    const data = new FormData();
    const file = document.getElementById("data_file").files[0];
    setFilename(file.name);
    data.append("file", file);
    data.append("type", file.type);
    fetch('/api/ipfs', {
        body: data,
        method: 'POST'
    })
    .then(res => res.json())
    .then(res => {console.log(res); setPhotoloading(false); setFilehash(res)})
    .catch(error => error.message)
}

const storeFile = async (e) => {
    setPhotoloading(true);
    props.guer.methods.doAddPhoto(props.mymia, filehash, filename, "image").send({
    from: props.myaccount
    })
    .then(function(result){
        setPhotoloading(false);
        window.location.reload(false);
    }).catch(function(error){
        console.log(error);
    });
    setPhotoloading(false);
}

const deleteFile = async (e) => {
    console.log(e);
    if (window.confirm("Do you want to delete this file?") === true){
        setPhotoloading(true);
        props.guer.methods.doRemovePhoto(props.mymia, e).send({
        from: props.myaccount
        })
        .then(function(result){
            console.log(result);
            setPhotoloading(false);
            window.location.reload(false);
        }).catch(function(error){
            console.log(error);
        });
        setPhotoloading(false);
    }
}

  const handleDropdownchange = (event) => {
    setFiletype(event.target.value);
  };

return (
    <Grid container spacing={1}>
    <Grid item xs={12}>
        <GuerNav image={props.image} address={props.mymia}/>
        <MinorHeader title={"Data"}/>
    </Grid>
        <Grid item xs={6} className="center">
            <div className="file-field input-field">
                <input
                    accept="image/*"
                    id="data_file"
                    multiple
                    type="file"
                    style={{display:'none'}}
                onChange={uploadFile}
                />
                <label htmlFor="data_file">
                <Button className="main-buttons" component="span">
                    <div>
                    <ArrowUpward/>
                    </div>
                    <div>Add Files</div>
                </Button>
                </label>
            </div>
        </Grid>
        <Grid item xs={6} className="center">
            <Button disabled={photoloading} className="main-buttons" type="submit" name="submit" id="submit" onClick={() => storeFile()}>{photoloading ? "Uploading..." : "Upload to Mia!"} </Button>
        </Grid>
        <Grid item xs={12}>
            <div style={{display:"flex",flexWrap:"wrap",textAlign:"center",margin:"0 auto"}}>
                <div style={{flex:'0 0 100%', height:'50px'}}>
                  <FormControl>
                    <InputLabel id="simple-select-label" >Filter</InputLabel>
                    <Select
                      labelId="simple-select-label"
                      id="simple-select"
                      value={filetype}
                      onChange={handleDropdownchange}
                    >
                      <MenuItem value={''}>All</MenuItem>
                      <MenuItem value={'image'}>Photos</MenuItem>
                      <MenuItem value={'videos'}>Videos</MenuItem>
                      <MenuItem value={'personal'}>Personal</MenuItem>
                      <MenuItem value={'excalidraw'}>Excalidraw</MenuItem>
                    </Select>
                  </FormControl>
                  </div>
                {
                    photos.map((result, index) => {
                        if (result[0] !== null && result[0] !== undefined && result[0] !== "" && (result[3] === filetype || filetype === '')) {
                            return (
                                <div className="imagebox" key={result[0]}>
                                    <div>
                                        <Cancel style={{opacity:".5",float:"right",cursor:"pointer"}} onClick={() => deleteFile(index)}/>
                                        <a rel="noopener noreferrer" target="_blank" href={"https://gateway.ipfs.io/ipfs/" + result[0]}><img alt="pic" src={"https://gateway.ipfs.io/ipfs/" + result[0]} style={{width:'80%', color:"white",float:"left",marginLeft:"10px"}}/></a>
                                        {/*
                                        <a rel="noopener noreferrer" target="_blank" href={"https://s3.cloudnado.com/guerbucket/" + result[0] + "?token=WV6KMH675PIK6ESVL2SWURLL"}><img alt="profile" src={"https://s3.cloudnado.com/guerbucket/" + result[0] + "?token=WV6KMH675PIK6ESVL2SWURLL"} style={{width:'80%'}}/></a>
                                        */}
                                    </div>
                                    <div>
                                        {result[1].substring(0,5)}...{result[1].substring(38,42)}
                                        <br/>
                                        <Moment format="MM/DD/YY HH:mm" tz="America/Los_Angeles" unix>{result[2]}</Moment>
                                    </div>
                                </div>
                                );
                        }
                        return('');
                    })
                }
            </div>      
        </Grid>      
    </Grid>
);
}

export default Photos;
