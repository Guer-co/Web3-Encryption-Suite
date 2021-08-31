import React, { useEffect} from 'react';
import PropTypes from 'prop-types';
import AppBar from '@material-ui/core/AppBar';
import Avatar from '@material-ui/core/Avatar';
import Drawer from '@material-ui/core/Drawer';
import IconButton from '@material-ui/core/IconButton';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';
import Grid from '@material-ui/core/Grid';
import CloseIcon from '@material-ui/icons/Close';
import Toolbar from '@material-ui/core/Toolbar';
import { makeStyles } from '@material-ui/core/styles';
import { Menu} from '@material-ui/icons';
import {Link } from 'react-router-dom';
import GuerLogo from '../images/logo.png';
//import MenuIcon from '@material-ui/icons/Menu';
//import Typography from '@material-ui/core/Typography';
//import Hidden from '@material-ui/core/Hidden';
//import CssBaseline from '@material-ui/core/CssBaseline';

const drawerWidth = 200;

const useStyles = makeStyles(theme => ({
root: {
    display: 'flex',
},
drawer: {
    [theme.breakpoints.up('sm')]: {
    width: drawerWidth,
    flexShrink: 0,
    },
},
appBar: {
    zIndex: theme.zIndex.drawer + 1,
},
menuButton: {
    marginRight: theme.spacing(2),
    [theme.breakpoints.up('sm')]: {
    display: 'none',
    },
},
toolbar: theme.mixins.toolbar,
drawerPaper: {
    width: drawerWidth,
    backgroundColor:"darkblue"
},
content: {
    flexGrow: 1,
    padding: theme.spacing(3),
},
closeMenuButton: {
    marginRight: 'auto',
    marginLeft: 0,
},
}));

const GuerNav = props => {

const classes = useStyles();
const [mobileOpen, setMobileOpen] = React.useState(false);

function handleDrawerToggle() {
    setMobileOpen(!mobileOpen)
}

useEffect(() => {
    
},[]);

const drawer = (
    <div>
    <List >
        <ListItem button>
            <Link to="/" style={{textDecoration:'none', color:'white'}}>
            <ListItemText primary="home" />
            </Link>
        </ListItem>
        <ListItem button>
            <Link to="/data" style={{textDecoration:'none', color:'white'}}>
            <ListItemText primary="my data" />
            </Link>
        </ListItem>
        <ListItem button>
            <Link to="/access" style={{textDecoration:'none', color:'white'}}>
            <ListItemText primary="my access" />
            </Link>
        </ListItem>
        <ListItem button>
            <Link to="/requests" style={{textDecoration:'none', color:'white'}}>
            <ListItemText primary="my access" />
            </Link>
        </ListItem>
        <ListItem button>
            <Link to="/access" style={{textDecoration:'none', color:'white'}}>
            <ListItemText primary="my requests" />
            </Link>
        </ListItem>
        <ListItem button>
            <Link to="/upload" style={{textDecoration:'none', color:'white'}}>
            <ListItemText primary="upload" />
            </Link>
        </ListItem>
        <ListItem button>
            <Link to="/profile" style={{textDecoration:'none', color:'white'}}>
            <ListItemText primary="profile" />
            </Link>
        </ListItem>
        <ListItem button>
            <Link to="/funds" style={{textDecoration:'none', color:'white'}}>
            <ListItemText primary="funds" />
            </Link>
        </ListItem>
        <ListItem button>
            <Link to="/advanced" style={{textDecoration:'none', color:'white'}}>
            <ListItemText primary="advanced" />
            </Link>
        </ListItem>
        <ListItem button>
            <Link to="/settings" style={{textDecoration:'none', color:'white'}}>
            <ListItemText primary="settings" />
            </Link>
        </ListItem>
    </List>
    </div>
);
return (
    <div>

    <AppBar position="static" className="appbar">
        <Toolbar>
            <Grid container spacing={1}>
            <Grid item xs={2}>
            <div style={{display:"block",alignItems: "center",justifyContent: "center", verticalAlign: "middle"}}>
                <IconButton
                    color="inherit"
                    aria-label="Open drawer"
                    edge="start"
                    onClick={handleDrawerToggle}
                    style={{marginTop:"8px"}}
                >
                <Menu />
                </IconButton>
            </div>
            </Grid>
            <Grid item xs={6}>
                <div style={{display:"flex",alignItems: "center",justifyContent: "center",width:"100%",marginTop:"5px"}}>
                <p>{props.address ? props.address.substring(0, 5) + "..." + props.address.substring(38, 42) : "?¿?¿?¿" }</p>
                </div>
            </Grid>
            <Grid item xs={4}>
                <Avatar style={{width:"70px",height:"70px",float:"right",backgroundColor:'white'}} alt="noImage" 
                src={props.image ? "https://cloudflare-ipfs.com/ipfs/" + props.image : GuerLogo} 
                /> 
            </Grid>
            </Grid>
        </Toolbar>
    </AppBar> 
      
    <nav className={classes.drawer}>
        <Drawer
            variant="persistent"
            anchor="left"
            open={mobileOpen}
            onClose={handleDrawerToggle}
            classes={{
            paper: classes.drawerPaper,
            }}
            ModalProps={{
            keepMounted: true,
            }}
        >
            <IconButton onClick={handleDrawerToggle} className={classes.closeMenuButton}>
            <CloseIcon style={{color:"white",paddingTop:"10px"}}/>
            </IconButton>
            {drawer}
        </Drawer>
    </nav>
    </div>
);
}
GuerNav.propTypes = {
// Injected by the documentation to work in an iframe.
// You won't need it on your project.
container: PropTypes.object,
};
export default GuerNav;