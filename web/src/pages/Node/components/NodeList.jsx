import { makeStyles } from '@material-ui/core'
import PropTypes from 'prop-types'
import React, {
  useEffect,
  useState
} from 'react'

import {
  Table, TableBody, TableCell, TableContainer, TableHead, TableRow, Paper,
  IconButton,
  Dialog,
  TextField
} from '@material-ui/core'
import DeleteIcon from '@material-ui/icons/Delete'
import PlusOneOutlinedIcon from '@material-ui/icons/PlusOneOutlined';
const useStyle = makeStyles((theme) => ({
  root: {
    minWidth: 650,
    '& .MuiIconButton-root': {
    }
  }
}))

const dialogStyle = makeStyles((theme) => ({
  root: {
    '& > *': {
      margin: theme.spacing(1)
    },
    '& .MuiDialog-paper': {
      padding: theme.spacing(1)
    },
    '& .MuiInput-root': {
      minWidth: 400,
    },
  }
}))

const AddProxyNodeDialog = (props) => {

  const defaultValue = ""

  const { open, onSubmit, onClose } = props

  const classes = dialogStyle()

  const [remote, setRemote] = useState(defaultValue)

  const handleClose = () => {
    // clean the form
    setRemote(defaultValue)
    onClose()
  }

  const handleChange = (e) => {
    setRemote(e.target.value)
  }

  const handleKeyPress = (e) => {
    if (e.code === 'Enter') {
      onSubmit(remote)
      setRemote(defaultValue)
    }
  }

  return (
    <Dialog className={classes.root} onClose={handleClose} aria-labelledby="add-proxy-node-dialog" open={open}>
      {/* <DialogTitle id="add-proxy-node-title">{"Add Proxy Node"}</DialogTitle> */}
      <form action="">
        <TextField
          id="remote-addr-input"
          required
          label="Remote Address"
          value={remote}
          onChange={handleChange}
          onKeyPress={handleKeyPress}
          placeholder={"http://127.0.0.1:8080"}
        />
      </form>
    </Dialog>
  )
}

AddProxyNodeDialog.propTypes = {
  open: PropTypes.bool.isRequired,
  onSubmit: PropTypes.func.isRequired,
  onClose: PropTypes.func.isRequired
}

const NodeList = (props) => {

  console.log(props)

  const classes = useStyle()

  const [open, setOpen] = useState(false)

  useEffect(() => {
    console.log(props)
    return () => {
    }
  }, [])

  const openAddNodePage = () => {
    setOpen(true)
  }

  const closeAddNodePage = () => {
    setOpen(false)
  }

  const addProxyNode = (remote) => {
    setOpen(false)
    props.addProxyNode(remote)
  }

  return (
    <TableContainer component={Paper}>
      <Table className={classes.root} aria-label="simple table">
        <TableHead>
          <TableRow>
            <TableCell>ID</TableCell>
            <TableCell align='left'>Remote Address</TableCell>
            <TableCell align='right'>
              <IconButton
                onClick={openAddNodePage}
              >
                <PlusOneOutlinedIcon
                  color='primary'
                />
              </IconButton>
            </TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {
            props.nodes.map((item) => {
              return <TableRow key={item.remote_addr}>
                <TableCell component="th" scope="row">{item.id}</TableCell>
                <TableCell align="left">{item.remote_addr}</TableCell>
                <TableCell align="right">
                  {/* Delete */}
                  <IconButton
                    aria-label="delete"
                    onClick={() => {
                      props.removeProxyNode(item.remote_addr)
                    }}
                  >
                    <DeleteIcon style={{
                      color: '#1C8185'
                    }} />
                  </IconButton>
                </TableCell>
              </TableRow>
            })
          }
        </TableBody>
      </Table>
      <AddProxyNodeDialog onClose={closeAddNodePage} open={open} onSubmit={addProxyNode} />
    </TableContainer>
  )
}

export {
  AddProxyNodeDialog,
  NodeList
}