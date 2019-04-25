CI/CD tools in which you can write your own pipelines in golang code (idea from Gaia where you can write it in any code of your choosing)
The pro of this is you are not subjected to anyones shitty languages for writing your pipelines (I'm looking at you jenkins)
The con of this is that you have to build your own libraries for all the different actions you might want to take


## MVP Goals
* Master/Minion Arch: One master only, Potentially many minions
* Master will round robin jobs to minions
* Pipeline code will only be able to be imported from github
* Use Hashicorp's plugin framework to allow pipelines to be compiled and ran
* The ability to trigger a recompile of the repo
* Built using all GRPC
* Masters on first creation of pipeline will grab from github and compile
* It will run under the user of the program and store the compiled code in a specific directory
* Maybe use boltdb so that the database is self contained?
* Pipelines can be created from different github branches
* Workers will connect with a shared secret for now
* workers will have a keep alive thread. The master will check the last checkin time of a minion before
    it round robins a request to it
* workers can be "suspended" which means jobs will no longer be sent to them
* The worker list and data is kept in memory by the master
* How do workers communicate they are not in a healthy enough state to continue operation?


## Far Far Future Goals
* Master/Minion Arch with Masters able to be distributed by paxos
* Authenication and ACLs can build upon github groups and oauth
* Auto recompilation of github repo

-----

Pipeline object {
    id
    name
    description
    last_run
    last_run_by
    created
    modified
    depends_on
}

worker object {
    id
    hostname
    state

}
