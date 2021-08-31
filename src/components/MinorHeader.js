import React from 'react';
import {Link } from 'react-router-dom';
import { ArrowBack } from '@material-ui/icons';

const MinorHeader = props =>  
    <div className="MinorHeader">
        <div style={{float:"left"}}><Link to="/" style={{color:"white"}}><ArrowBack style={{fontSize:"40px"}}/></Link></div>
        <h3 style={{paddingTop:'5px'}}>&nbsp;{props.title}</h3>
    </div>
;

export default MinorHeader;