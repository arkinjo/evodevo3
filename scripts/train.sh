#!/bin/zsh

export GOMAXPROCS=8

BINDIR=~/work/GitHub/evodevo3/bin
genenv=${BINDIR}/genenv
runsim=${BINDIR}/runsim
gpplot=${BINDIR}/gpplot

POPSIZE=500
ENVSFILE=Environments.json
${genenv} -o ${ENVSFILE} -n 50 -denv 0.5 -seed 13 2> environ.log

EI1=30
EI0=`printf "%2.2d" $[${EI1}-1]`
EI2=50
EEND=`printf "%2.2d" $[${EI2}-1]`
MODELS=(Full)
#MODELS=(Full NoCue NoDev)
#MODELS=(Hie0 Hie1 Hie2)
#MODELS=(NullCue NullDev NullHie)
for model in $MODELS; do
    ${runsim} -envs ${ENVSFILE} -model ${model} -popsize ${POPSIZE} \
	      -env_start 0 -env_end ${EI1} -ngen 200 \
	> data/${model}_train.out
done 2> train.log

