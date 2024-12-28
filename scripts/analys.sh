#!/bin/zsh

export GOMAXPROCS=8

BINDIR=~/work/GitHub/evodevo3/bin
runsim=${BINDIR}/runsim
gpplot=${BINDIR}/gpplot

POPSIZE=500
ENVSFILE=Environments.json
EI1=30
EI0=`printf "%2.2d" $[${EI1}-1]`
EI2=31
EEND=`printf "%2.2d" $[${EI2}-1]`
#MODELS=(Full Null NoCue NoDev Hie0 Hie1 Hie2 NullHie NullCue NullDev)
#MODELS=(Full Null NullCue NoCue NoDev Hie0 Hie1 Hie2)
MODELS=(NullCue)
for model in $MODELS; do
    for epoch in {${EI1}..${EEND}}; do
	echo epoch $epoch
	${gpplot} -setting traj/Setting_${model}.json \
 		  -envs ${ENVSFILE} -ienv ${epoch} traj/${model}_${epoch}_*.traj.gz
    done 2> plot.log
done
