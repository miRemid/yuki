import { makeStyles } from '@material-ui/core'
import React, { useEffect, useState } from 'react'

import {
  Paper,
  TableContainer, Table, TableHead, TableRow, TableCell, TableBody,
  IconButton, Button,
  Dialog,
  Select, MenuItem, InputLabel, Input
} from '@material-ui/core'

import PlusOneOutlinedIcon from '@material-ui/icons/PlusOneOutlined'
import DeleteIcon from '@material-ui/icons/Delete'
import BackupIcon from '@material-ui/icons/Backup'

const useStyle = makeStyles((theme) => ({
  root: {
    minWidth: 650,
  }
}))

const dialogStyle = makeStyles((theme) => ({
  root: {
    '& > *': {
      margin: theme.spacing(1)
    },
    '& .MuiDialog-paper': {
      padding: theme.spacing(1),
      minWidth: 400,
    },
    '& .input': {
      marginBottom: '1%',
      '& .MuiInput-root': {
        width: '100%'
      }
    },
    '& .form': {
      display: 'flex',
      flexDirection: 'column',
    },
    '& .MuiButton-root': {
      marginTop: theme.spacing(1),
    }
  }
}))

const AddRuleDialog = (props) => {

  const classes = dialogStyle()

  const { open, onSubmit, onClose, nodes } = props

  const [cmd, setCmd] = useState('')
  const [remote, setRemtoe] = useState('')

  const handleCmdChange = (e) => {
    setCmd(e.target.value)
  }

  const handleRemoteChange = (e) => {
    setRemtoe(e.target.value)
  }

  const handleClose = () => {
    setCmd('')
    setRemtoe('')
    onClose()
  }

  const handleUpload = () => {
    console.log(cmd, remote)
    onSubmit(cmd, remote)
    handleClose()
  }

  return (
    <Dialog className={classes.root} onClose={handleClose} aira-labelledby="add-rule-dialog" open={open}>
      <form className={"form"}>
        <div className="input">
          <InputLabel id="demo-simple-select-label">Command</InputLabel>
          <Input
            autoFocus
            required
            value={cmd}
            onChange={handleCmdChange}
          />
        </div>
        <div className="input">
          <InputLabel id="demo-simple-select-label">Remote Address</InputLabel>
          <Select
            labelId="demo-simple-select-label"
            id="demo-simple-select"
            value={remote}
            onChange={handleRemoteChange}
          >
            {
              nodes.map((item) => {
                return <MenuItem key={item.remote_addr} value={item.remote_addr}>{item.remote_addr}</MenuItem>
              })
            }
          </Select>
        </div>

        <Button
          variant="contained"
          color="primary"
          className={classes.button}
          startIcon={<BackupIcon />}
          onClick={handleUpload}
        >
        </Button>
      </form>
    </Dialog>
  )
}

AddRuleDialog.propTypes = {

}

const RuleList = (props) => {
  const classes = useStyle()

  const [open, setOpen] = useState(false)

  useEffect(() => {

  }, [])

  const openAddRuleDialog = () => {
    setOpen(true)
  }

  const closeAddRuleDialog = () => {
    setOpen(false)
  }

  return (
    <TableContainer component={Paper}>
      <Table className={classes.root} aria-label="rule table">
        <TableHead>
          <TableRow>
            <TableCell align='left'>Command</TableCell>
            <TableCell align='center'>Remote Address</TableCell>
            <TableCell align='right'>
              <IconButton
                onClick={openAddRuleDialog}
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
            props.rules.map((item) => {
              return <TableRow key={item.remote_addr}>
                {/* <TableCell component="th" scope="row">{item.id}</TableCell> */}
                <TableCell align="left">{item.cmd}</TableCell>
                <TableCell align="center">{item.remote_addr}</TableCell>
                <TableCell align="right">
                  <IconButton
                    aria-label='delete'
                    onClick={() => {
                      props.removeRule(item.cmd)
                    }}
                  >
                    <DeleteIcon
                      style={{
                        color: '#1c8185'
                      }}
                    />
                  </IconButton>
                </TableCell>
              </TableRow>
            })
          }
        </TableBody>
      </Table>
      <AddRuleDialog
        onClose={closeAddRuleDialog}
        open={open}
        onSubmit={props.addRule}
        nodes={props.nodes}
      ></AddRuleDialog>
    </TableContainer>
  )

}

export {
  AddRuleDialog,
  RuleList
}