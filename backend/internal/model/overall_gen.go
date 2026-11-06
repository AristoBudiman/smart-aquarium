// Kode dibuat oleh ml/train_export.py. Jangan diedit.
package model

func ScoreOverall(input []float64) []float64 {
    var var0 []float64
    if input[1] <= 24.975000381469727 {
        var0 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[0] <= 32.010000228881836 {
            if input[1] <= 50.27000045776367 {
                var0 = []float64{0.0, 1.0, 0.0, 0.0}
            } else {
                if input[1] <= 151.63999938964844 {
                    var0 = []float64{0.0, 0.0, 1.0, 0.0}
                } else {
                    var0 = []float64{0.0, 0.0, 0.0, 1.0}
                }
            }
        } else {
            if input[1] <= 150.37999725341797 {
                var0 = []float64{0.0, 0.0, 1.0, 0.0}
            } else {
                var0 = []float64{0.0, 0.0, 0.0, 1.0}
            }
        }
    }
    var var1 []float64
    if input[0] <= 23.010000228881836 {
        if input[0] <= 20.015000343322754 {
            var1 = []float64{0.0, 0.0, 0.0, 1.0}
        } else {
            var1 = []float64{0.0, 0.0, 1.0, 0.0}
        }
    } else {
        if input[0] <= 32.010000228881836 {
            if input[0] <= 25.0 {
                var1 = []float64{0.0, 1.0, 0.0, 0.0}
            } else {
                if input[1] <= 25.0 {
                    var1 = []float64{1.0, 0.0, 0.0, 0.0}
                } else {
                    var1 = []float64{0.0, 1.0, 0.0, 0.0}
                }
            }
        } else {
            if input[1] <= 150.37999725341797 {
                var1 = []float64{0.0, 0.0, 1.0, 0.0}
            } else {
                var1 = []float64{0.0, 0.0, 0.0, 1.0}
            }
        }
    }
    var var2 []float64
    if input[0] <= 22.994999885559082 {
        if input[1] <= 150.56999969482422 {
            var2 = []float64{0.0, 0.0, 1.0, 0.0}
        } else {
            var2 = []float64{0.0, 0.0, 0.0, 1.0}
        }
    } else {
        if input[1] <= 24.984999656677246 {
            var2 = []float64{1.0, 0.0, 0.0, 0.0}
        } else {
            if input[1] <= 50.09499931335449 {
                var2 = []float64{0.0, 1.0, 0.0, 0.0}
            } else {
                if input[0] <= 33.9950008392334 {
                    var2 = []float64{0.0, 0.0, 1.0, 0.0}
                } else {
                    var2 = []float64{0.0, 0.0, 0.0, 1.0}
                }
            }
        }
    }
    var var3 []float64
    if input[0] <= 34.0049991607666 {
        if input[0] <= 32.010000228881836 {
            if input[1] <= 24.984999656677246 {
                var3 = []float64{1.0, 0.0, 0.0, 0.0}
            } else {
                if input[1] <= 50.204999923706055 {
                    var3 = []float64{0.0, 1.0, 0.0, 0.0}
                } else {
                    if input[0] <= 20.074999809265137 {
                        var3 = []float64{0.0, 0.0, 0.0, 1.0}
                    } else {
                        var3 = []float64{0.0, 0.0, 1.0, 0.0}
                    }
                }
            }
        } else {
            var3 = []float64{0.0, 0.0, 1.0, 0.0}
        }
    } else {
        var3 = []float64{0.0, 0.0, 0.0, 1.0}
    }
    var var4 []float64
    if input[1] <= 24.985000610351562 {
        var4 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[1] <= 50.15999984741211 {
            var4 = []float64{0.0, 1.0, 0.0, 0.0}
        } else {
            if input[1] <= 150.68000030517578 {
                var4 = []float64{0.0, 0.0, 1.0, 0.0}
            } else {
                var4 = []float64{0.0, 0.0, 0.0, 1.0}
            }
        }
    }
    var var5 []float64
    if input[1] <= 24.994999885559082 {
        var5 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[0] <= 20.015000343322754 {
            var5 = []float64{0.0, 0.0, 0.0, 1.0}
        } else {
            if input[1] <= 50.18499946594238 {
                var5 = []float64{0.0, 1.0, 0.0, 0.0}
            } else {
                if input[1] <= 152.40499877929688 {
                    var5 = []float64{0.0, 0.0, 1.0, 0.0}
                } else {
                    var5 = []float64{0.0, 0.0, 0.0, 1.0}
                }
            }
        }
    }
    var var6 []float64
    if input[0] <= 34.0049991607666 {
        if input[0] <= 22.994999885559082 {
            if input[0] <= 19.954999923706055 {
                var6 = []float64{0.0, 0.0, 0.0, 1.0}
            } else {
                var6 = []float64{0.0, 0.0, 1.0, 0.0}
            }
        } else {
            if input[1] <= 24.984999656677246 {
                var6 = []float64{1.0, 0.0, 0.0, 0.0}
            } else {
                if input[0] <= 31.994999885559082 {
                    var6 = []float64{0.0, 1.0, 0.0, 0.0}
                } else {
                    var6 = []float64{0.0, 0.0, 1.0, 0.0}
                }
            }
        }
    } else {
        var6 = []float64{0.0, 0.0, 0.0, 1.0}
    }
    var var7 []float64
    if input[0] <= 32.00999927520752 {
        if input[1] <= 24.984999656677246 {
            var7 = []float64{1.0, 0.0, 0.0, 0.0}
        } else {
            if input[0] <= 22.994999885559082 {
                if input[1] <= 150.56999969482422 {
                    var7 = []float64{0.0, 0.0, 1.0, 0.0}
                } else {
                    var7 = []float64{0.0, 0.0, 0.0, 1.0}
                }
            } else {
                var7 = []float64{0.0, 1.0, 0.0, 0.0}
            }
        }
    } else {
        if input[1] <= 150.37999725341797 {
            var7 = []float64{0.0, 0.0, 1.0, 0.0}
        } else {
            var7 = []float64{0.0, 0.0, 0.0, 1.0}
        }
    }
    var var8 []float64
    if input[0] <= 22.994999885559082 {
        if input[0] <= 20.010000228881836 {
            var8 = []float64{0.0, 0.0, 0.0, 1.0}
        } else {
            var8 = []float64{0.0, 0.0, 1.0, 0.0}
        }
    } else {
        if input[0] <= 32.010000228881836 {
            if input[1] <= 24.984999656677246 {
                var8 = []float64{1.0, 0.0, 0.0, 0.0}
            } else {
                var8 = []float64{0.0, 1.0, 0.0, 0.0}
            }
        } else {
            if input[0] <= 34.0 {
                var8 = []float64{0.0, 0.0, 1.0, 0.0}
            } else {
                var8 = []float64{0.0, 0.0, 0.0, 1.0}
            }
        }
    }
    var var9 []float64
    if input[1] <= 24.94499969482422 {
        var9 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[0] <= 33.9950008392334 {
            if input[1] <= 50.10000038146973 {
                var9 = []float64{0.0, 1.0, 0.0, 0.0}
            } else {
                if input[1] <= 151.9800033569336 {
                    var9 = []float64{0.0, 0.0, 1.0, 0.0}
                } else {
                    var9 = []float64{0.0, 0.0, 0.0, 1.0}
                }
            }
        } else {
            var9 = []float64{0.0, 0.0, 0.0, 1.0}
        }
    }
    var var10 []float64
    if input[0] <= 34.0 {
        if input[1] <= 24.99000072479248 {
            var10 = []float64{1.0, 0.0, 0.0, 0.0}
        } else {
            if input[1] <= 50.119998931884766 {
                var10 = []float64{0.0, 1.0, 0.0, 0.0}
            } else {
                if input[1] <= 150.95999908447266 {
                    var10 = []float64{0.0, 0.0, 1.0, 0.0}
                } else {
                    var10 = []float64{0.0, 0.0, 0.0, 1.0}
                }
            }
        }
    } else {
        var10 = []float64{0.0, 0.0, 0.0, 1.0}
    }
    var var11 []float64
    if input[1] <= 24.984999656677246 {
        var11 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[0] <= 33.98499870300293 {
            if input[1] <= 50.18499946594238 {
                var11 = []float64{0.0, 1.0, 0.0, 0.0}
            } else {
                if input[0] <= 20.074999809265137 {
                    var11 = []float64{0.0, 0.0, 0.0, 1.0}
                } else {
                    var11 = []float64{0.0, 0.0, 1.0, 0.0}
                }
            }
        } else {
            var11 = []float64{0.0, 0.0, 0.0, 1.0}
        }
    }
    var var12 []float64
    if input[0] <= 32.010000228881836 {
        if input[1] <= 24.994999885559082 {
            var12 = []float64{1.0, 0.0, 0.0, 0.0}
        } else {
            if input[1] <= 50.204999923706055 {
                var12 = []float64{0.0, 1.0, 0.0, 0.0}
            } else {
                if input[1] <= 150.95999908447266 {
                    var12 = []float64{0.0, 0.0, 1.0, 0.0}
                } else {
                    var12 = []float64{0.0, 0.0, 0.0, 1.0}
                }
            }
        }
    } else {
        if input[1] <= 152.70499420166016 {
            var12 = []float64{0.0, 0.0, 1.0, 0.0}
        } else {
            var12 = []float64{0.0, 0.0, 0.0, 1.0}
        }
    }
    var var13 []float64
    if input[0] <= 22.994999885559082 {
        if input[1] <= 150.6199951171875 {
            var13 = []float64{0.0, 0.0, 1.0, 0.0}
        } else {
            var13 = []float64{0.0, 0.0, 0.0, 1.0}
        }
    } else {
        if input[1] <= 24.994999885559082 {
            var13 = []float64{1.0, 0.0, 0.0, 0.0}
        } else {
            if input[0] <= 32.045000076293945 {
                if input[1] <= 81.67500114440918 {
                    var13 = []float64{0.0, 1.0, 0.0, 0.0}
                } else {
                    var13 = []float64{0.0, 0.0, 1.0, 0.0}
                }
            } else {
                if input[0] <= 34.01499938964844 {
                    if input[0] <= 33.96500015258789 {
                        var13 = []float64{0.0, 0.0, 1.0, 0.0}
                    } else {
                        if input[1] <= 344.1600036621094 {
                            var13 = []float64{0.0, 0.0, 1.0, 0.0}
                        } else {
                            var13 = []float64{0.0, 0.0, 0.0, 1.0}
                        }
                    }
                } else {
                    var13 = []float64{0.0, 0.0, 0.0, 1.0}
                }
            }
        }
    }
    var var14 []float64
    if input[1] <= 24.975000381469727 {
        var14 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[1] <= 50.119998931884766 {
            var14 = []float64{0.0, 1.0, 0.0, 0.0}
        } else {
            if input[1] <= 150.68000030517578 {
                var14 = []float64{0.0, 0.0, 1.0, 0.0}
            } else {
                var14 = []float64{0.0, 0.0, 0.0, 1.0}
            }
        }
    }
    var var15 []float64
    if input[0] <= 22.994999885559082 {
        if input[0] <= 20.074999809265137 {
            var15 = []float64{0.0, 0.0, 0.0, 1.0}
        } else {
            var15 = []float64{0.0, 0.0, 1.0, 0.0}
        }
    } else {
        if input[0] <= 31.994999885559082 {
            if input[1] <= 24.984999656677246 {
                var15 = []float64{1.0, 0.0, 0.0, 0.0}
            } else {
                var15 = []float64{0.0, 1.0, 0.0, 0.0}
            }
        } else {
            if input[0] <= 33.9950008392334 {
                if input[1] <= 40.59499931335449 {
                    var15 = []float64{0.0, 1.0, 0.0, 0.0}
                } else {
                    var15 = []float64{0.0, 0.0, 1.0, 0.0}
                }
            } else {
                if input[1] <= 136.31999969482422 {
                    var15 = []float64{0.0, 0.0, 1.0, 0.0}
                } else {
                    var15 = []float64{0.0, 0.0, 0.0, 1.0}
                }
            }
        }
    }
    var var16 []float64
    if input[1] <= 24.975000381469727 {
        var16 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[1] <= 50.09499931335449 {
            var16 = []float64{0.0, 1.0, 0.0, 0.0}
        } else {
            if input[1] <= 151.68000030517578 {
                var16 = []float64{0.0, 0.0, 1.0, 0.0}
            } else {
                var16 = []float64{0.0, 0.0, 0.0, 1.0}
            }
        }
    }
    var var17 []float64
    if input[0] <= 22.994999885559082 {
        if input[0] <= 20.074999809265137 {
            var17 = []float64{0.0, 0.0, 0.0, 1.0}
        } else {
            var17 = []float64{0.0, 0.0, 1.0, 0.0}
        }
    } else {
        if input[0] <= 32.01499938964844 {
            if input[1] <= 24.994999885559082 {
                var17 = []float64{1.0, 0.0, 0.0, 0.0}
            } else {
                if input[1] <= 81.67500114440918 {
                    var17 = []float64{0.0, 1.0, 0.0, 0.0}
                } else {
                    var17 = []float64{0.0, 0.0, 1.0, 0.0}
                }
            }
        } else {
            if input[0] <= 33.9950008392334 {
                var17 = []float64{0.0, 0.0, 1.0, 0.0}
            } else {
                var17 = []float64{0.0, 0.0, 0.0, 1.0}
            }
        }
    }
    var var18 []float64
    if input[0] <= 22.994999885559082 {
        if input[0] <= 20.015000343322754 {
            var18 = []float64{0.0, 0.0, 0.0, 1.0}
        } else {
            var18 = []float64{0.0, 0.0, 1.0, 0.0}
        }
    } else {
        if input[1] <= 24.984999656677246 {
            var18 = []float64{1.0, 0.0, 0.0, 0.0}
        } else {
            if input[0] <= 32.010000228881836 {
                var18 = []float64{0.0, 1.0, 0.0, 0.0}
            } else {
                if input[0] <= 34.0049991607666 {
                    var18 = []float64{0.0, 0.0, 1.0, 0.0}
                } else {
                    var18 = []float64{0.0, 0.0, 0.0, 1.0}
                }
            }
        }
    }
    var var19 []float64
    if input[0] <= 34.010000228881836 {
        if input[0] <= 22.994999885559082 {
            if input[1] <= 150.6199951171875 {
                var19 = []float64{0.0, 0.0, 1.0, 0.0}
            } else {
                var19 = []float64{0.0, 0.0, 0.0, 1.0}
            }
        } else {
            if input[0] <= 31.994999885559082 {
                if input[1] <= 24.994999885559082 {
                    var19 = []float64{1.0, 0.0, 0.0, 0.0}
                } else {
                    var19 = []float64{0.0, 1.0, 0.0, 0.0}
                }
            } else {
                var19 = []float64{0.0, 0.0, 1.0, 0.0}
            }
        }
    } else {
        var19 = []float64{0.0, 0.0, 0.0, 1.0}
    }
    var var20 []float64
    if input[1] <= 24.984999656677246 {
        var20 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[0] <= 20.015000343322754 {
            var20 = []float64{0.0, 0.0, 0.0, 1.0}
        } else {
            if input[0] <= 32.01499938964844 {
                if input[0] <= 22.994999885559082 {
                    var20 = []float64{0.0, 0.0, 1.0, 0.0}
                } else {
                    if input[0] <= 31.994999885559082 {
                        var20 = []float64{0.0, 1.0, 0.0, 0.0}
                    } else {
                        if input[1] <= 72.1500015258789 {
                            var20 = []float64{0.0, 1.0, 0.0, 0.0}
                        } else {
                            var20 = []float64{0.0, 0.0, 1.0, 0.0}
                        }
                    }
                }
            } else {
                if input[1] <= 150.37999725341797 {
                    var20 = []float64{0.0, 0.0, 1.0, 0.0}
                } else {
                    var20 = []float64{0.0, 0.0, 0.0, 1.0}
                }
            }
        }
    }
    var var21 []float64
    if input[1] <= 24.984999656677246 {
        var21 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[0] <= 34.0 {
            if input[0] <= 22.994999885559082 {
                if input[1] <= 150.95999908447266 {
                    var21 = []float64{0.0, 0.0, 1.0, 0.0}
                } else {
                    var21 = []float64{0.0, 0.0, 0.0, 1.0}
                }
            } else {
                if input[1] <= 50.18499946594238 {
                    var21 = []float64{0.0, 1.0, 0.0, 0.0}
                } else {
                    var21 = []float64{0.0, 0.0, 1.0, 0.0}
                }
            }
        } else {
            var21 = []float64{0.0, 0.0, 0.0, 1.0}
        }
    }
    var var22 []float64
    if input[0] <= 22.994999885559082 {
        if input[1] <= 150.95999908447266 {
            var22 = []float64{0.0, 0.0, 1.0, 0.0}
        } else {
            var22 = []float64{0.0, 0.0, 0.0, 1.0}
        }
    } else {
        if input[0] <= 31.994999885559082 {
            if input[0] <= 25.014999389648438 {
                var22 = []float64{0.0, 1.0, 0.0, 0.0}
            } else {
                if input[1] <= 24.985000610351562 {
                    var22 = []float64{1.0, 0.0, 0.0, 0.0}
                } else {
                    var22 = []float64{0.0, 1.0, 0.0, 0.0}
                }
            }
        } else {
            if input[0] <= 34.0 {
                if input[1] <= 40.65999984741211 {
                    var22 = []float64{0.0, 1.0, 0.0, 0.0}
                } else {
                    var22 = []float64{0.0, 0.0, 1.0, 0.0}
                }
            } else {
                var22 = []float64{0.0, 0.0, 0.0, 1.0}
            }
        }
    }
    var var23 []float64
    if input[1] <= 24.994999885559082 {
        var23 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[1] <= 50.119998931884766 {
            var23 = []float64{0.0, 1.0, 0.0, 0.0}
        } else {
            if input[1] <= 151.9800033569336 {
                var23 = []float64{0.0, 0.0, 1.0, 0.0}
            } else {
                var23 = []float64{0.0, 0.0, 0.0, 1.0}
            }
        }
    }
    var var24 []float64
    if input[1] <= 24.994999885559082 {
        var24 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[0] <= 33.96999931335449 {
            if input[1] <= 50.15999984741211 {
                var24 = []float64{0.0, 1.0, 0.0, 0.0}
            } else {
                if input[0] <= 20.074999809265137 {
                    var24 = []float64{0.0, 0.0, 0.0, 1.0}
                } else {
                    var24 = []float64{0.0, 0.0, 1.0, 0.0}
                }
            }
        } else {
            if input[1] <= 136.31999969482422 {
                var24 = []float64{0.0, 0.0, 1.0, 0.0}
            } else {
                var24 = []float64{0.0, 0.0, 0.0, 1.0}
            }
        }
    }
    var var25 []float64
    if input[0] <= 33.9950008392334 {
        if input[0] <= 23.0 {
            if input[0] <= 20.074999809265137 {
                var25 = []float64{0.0, 0.0, 0.0, 1.0}
            } else {
                var25 = []float64{0.0, 0.0, 1.0, 0.0}
            }
        } else {
            if input[1] <= 24.994999885559082 {
                var25 = []float64{1.0, 0.0, 0.0, 0.0}
            } else {
                if input[1] <= 50.15999984741211 {
                    var25 = []float64{0.0, 1.0, 0.0, 0.0}
                } else {
                    var25 = []float64{0.0, 0.0, 1.0, 0.0}
                }
            }
        }
    } else {
        var25 = []float64{0.0, 0.0, 0.0, 1.0}
    }
    var var26 []float64
    if input[1] <= 24.984999656677246 {
        var26 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[1] <= 50.2450008392334 {
            var26 = []float64{0.0, 1.0, 0.0, 0.0}
        } else {
            if input[1] <= 150.95999908447266 {
                var26 = []float64{0.0, 0.0, 1.0, 0.0}
            } else {
                var26 = []float64{0.0, 0.0, 0.0, 1.0}
            }
        }
    }
    var var27 []float64
    if input[0] <= 22.994999885559082 {
        if input[0] <= 20.015000343322754 {
            var27 = []float64{0.0, 0.0, 0.0, 1.0}
        } else {
            var27 = []float64{0.0, 0.0, 1.0, 0.0}
        }
    } else {
        if input[0] <= 32.010000228881836 {
            if input[1] <= 24.984999656677246 {
                var27 = []float64{1.0, 0.0, 0.0, 0.0}
            } else {
                if input[0] <= 31.994999885559082 {
                    var27 = []float64{0.0, 1.0, 0.0, 0.0}
                } else {
                    if input[1] <= 69.7750015258789 {
                        var27 = []float64{0.0, 1.0, 0.0, 0.0}
                    } else {
                        var27 = []float64{0.0, 0.0, 1.0, 0.0}
                    }
                }
            }
        } else {
            if input[0] <= 33.9950008392334 {
                var27 = []float64{0.0, 0.0, 1.0, 0.0}
            } else {
                var27 = []float64{0.0, 0.0, 0.0, 1.0}
            }
        }
    }
    var var28 []float64
    if input[1] <= 24.984999656677246 {
        var28 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[1] <= 50.09499931335449 {
            var28 = []float64{0.0, 1.0, 0.0, 0.0}
        } else {
            if input[1] <= 150.6199951171875 {
                var28 = []float64{0.0, 0.0, 1.0, 0.0}
            } else {
                var28 = []float64{0.0, 0.0, 0.0, 1.0}
            }
        }
    }
    var var29 []float64
    if input[0] <= 22.994999885559082 {
        if input[0] <= 20.015000343322754 {
            var29 = []float64{0.0, 0.0, 0.0, 1.0}
        } else {
            var29 = []float64{0.0, 0.0, 1.0, 0.0}
        }
    } else {
        if input[0] <= 32.010000228881836 {
            if input[1] <= 24.984999656677246 {
                var29 = []float64{1.0, 0.0, 0.0, 0.0}
            } else {
                if input[1] <= 81.67500114440918 {
                    var29 = []float64{0.0, 1.0, 0.0, 0.0}
                } else {
                    var29 = []float64{0.0, 0.0, 1.0, 0.0}
                }
            }
        } else {
            if input[0] <= 33.9950008392334 {
                var29 = []float64{0.0, 0.0, 1.0, 0.0}
            } else {
                if input[1] <= 138.0449981689453 {
                    var29 = []float64{0.0, 0.0, 1.0, 0.0}
                } else {
                    var29 = []float64{0.0, 0.0, 0.0, 1.0}
                }
            }
        }
    }
    var var30 []float64
    if input[0] <= 34.0049991607666 {
        if input[1] <= 24.994999885559082 {
            var30 = []float64{1.0, 0.0, 0.0, 0.0}
        } else {
            if input[0] <= 22.994999885559082 {
                if input[0] <= 20.085000038146973 {
                    var30 = []float64{0.0, 0.0, 0.0, 1.0}
                } else {
                    var30 = []float64{0.0, 0.0, 1.0, 0.0}
                }
            } else {
                if input[1] <= 50.119998931884766 {
                    var30 = []float64{0.0, 1.0, 0.0, 0.0}
                } else {
                    if input[1] <= 358.2200012207031 {
                        var30 = []float64{0.0, 0.0, 1.0, 0.0}
                    } else {
                        var30 = []float64{0.0, 0.0, 0.0, 1.0}
                    }
                }
            }
        }
    } else {
        var30 = []float64{0.0, 0.0, 0.0, 1.0}
    }
    var var31 []float64
    if input[1] <= 24.994999885559082 {
        var31 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[0] <= 23.0 {
            if input[0] <= 20.074999809265137 {
                var31 = []float64{0.0, 0.0, 0.0, 1.0}
            } else {
                var31 = []float64{0.0, 0.0, 1.0, 0.0}
            }
        } else {
            if input[1] <= 50.32499885559082 {
                var31 = []float64{0.0, 1.0, 0.0, 0.0}
            } else {
                if input[0] <= 34.01499938964844 {
                    var31 = []float64{0.0, 0.0, 1.0, 0.0}
                } else {
                    var31 = []float64{0.0, 0.0, 0.0, 1.0}
                }
            }
        }
    }
    var var32 []float64
    if input[1] <= 24.984999656677246 {
        var32 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[0] <= 20.015000343322754 {
            var32 = []float64{0.0, 0.0, 0.0, 1.0}
        } else {
            if input[1] <= 50.18499946594238 {
                var32 = []float64{0.0, 1.0, 0.0, 0.0}
            } else {
                if input[1] <= 150.68000030517578 {
                    var32 = []float64{0.0, 0.0, 1.0, 0.0}
                } else {
                    var32 = []float64{0.0, 0.0, 0.0, 1.0}
                }
            }
        }
    }
    var var33 []float64
    if input[0] <= 22.994999885559082 {
        if input[1] <= 150.95999908447266 {
            var33 = []float64{0.0, 0.0, 1.0, 0.0}
        } else {
            var33 = []float64{0.0, 0.0, 0.0, 1.0}
        }
    } else {
        if input[1] <= 24.994999885559082 {
            var33 = []float64{1.0, 0.0, 0.0, 0.0}
        } else {
            if input[0] <= 31.994999885559082 {
                var33 = []float64{0.0, 1.0, 0.0, 0.0}
            } else {
                if input[0] <= 33.9950008392334 {
                    var33 = []float64{0.0, 0.0, 1.0, 0.0}
                } else {
                    if input[0] <= 34.0049991607666 {
                        if input[1] <= 344.1600036621094 {
                            var33 = []float64{0.0, 0.0, 1.0, 0.0}
                        } else {
                            var33 = []float64{0.0, 0.0, 0.0, 1.0}
                        }
                    } else {
                        var33 = []float64{0.0, 0.0, 0.0, 1.0}
                    }
                }
            }
        }
    }
    var var34 []float64
    if input[0] <= 22.994999885559082 {
        if input[0] <= 20.074999809265137 {
            var34 = []float64{0.0, 0.0, 0.0, 1.0}
        } else {
            var34 = []float64{0.0, 0.0, 1.0, 0.0}
        }
    } else {
        if input[1] <= 24.994999885559082 {
            var34 = []float64{1.0, 0.0, 0.0, 0.0}
        } else {
            if input[0] <= 32.010000228881836 {
                var34 = []float64{0.0, 1.0, 0.0, 0.0}
            } else {
                if input[1] <= 152.70499420166016 {
                    var34 = []float64{0.0, 0.0, 1.0, 0.0}
                } else {
                    var34 = []float64{0.0, 0.0, 0.0, 1.0}
                }
            }
        }
    }
    var var35 []float64
    if input[1] <= 24.994999885559082 {
        var35 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[1] <= 50.18000030517578 {
            var35 = []float64{0.0, 1.0, 0.0, 0.0}
        } else {
            if input[0] <= 33.9950008392334 {
                if input[1] <= 150.95999908447266 {
                    var35 = []float64{0.0, 0.0, 1.0, 0.0}
                } else {
                    var35 = []float64{0.0, 0.0, 0.0, 1.0}
                }
            } else {
                var35 = []float64{0.0, 0.0, 0.0, 1.0}
            }
        }
    }
    var var36 []float64
    if input[0] <= 32.01499938964844 {
        if input[1] <= 24.960000038146973 {
            var36 = []float64{1.0, 0.0, 0.0, 0.0}
        } else {
            if input[0] <= 22.994999885559082 {
                if input[0] <= 20.015000343322754 {
                    var36 = []float64{0.0, 0.0, 0.0, 1.0}
                } else {
                    var36 = []float64{0.0, 0.0, 1.0, 0.0}
                }
            } else {
                var36 = []float64{0.0, 1.0, 0.0, 0.0}
            }
        }
    } else {
        if input[1] <= 152.10499572753906 {
            var36 = []float64{0.0, 0.0, 1.0, 0.0}
        } else {
            var36 = []float64{0.0, 0.0, 0.0, 1.0}
        }
    }
    var var37 []float64
    if input[1] <= 24.994999885559082 {
        var37 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[1] <= 50.09499931335449 {
            var37 = []float64{0.0, 1.0, 0.0, 0.0}
        } else {
            if input[1] <= 150.65999603271484 {
                var37 = []float64{0.0, 0.0, 1.0, 0.0}
            } else {
                var37 = []float64{0.0, 0.0, 0.0, 1.0}
            }
        }
    }
    var var38 []float64
    if input[0] <= 22.994999885559082 {
        if input[1] <= 150.56999969482422 {
            var38 = []float64{0.0, 0.0, 1.0, 0.0}
        } else {
            var38 = []float64{0.0, 0.0, 0.0, 1.0}
        }
    } else {
        if input[0] <= 32.00500011444092 {
            if input[1] <= 24.994999885559082 {
                var38 = []float64{1.0, 0.0, 0.0, 0.0}
            } else {
                var38 = []float64{0.0, 1.0, 0.0, 0.0}
            }
        } else {
            if input[1] <= 152.04000091552734 {
                var38 = []float64{0.0, 0.0, 1.0, 0.0}
            } else {
                var38 = []float64{0.0, 0.0, 0.0, 1.0}
            }
        }
    }
    var var39 []float64
    if input[0] <= 22.994999885559082 {
        if input[1] <= 150.95999908447266 {
            var39 = []float64{0.0, 0.0, 1.0, 0.0}
        } else {
            var39 = []float64{0.0, 0.0, 0.0, 1.0}
        }
    } else {
        if input[1] <= 24.994999885559082 {
            var39 = []float64{1.0, 0.0, 0.0, 0.0}
        } else {
            if input[0] <= 32.01499938964844 {
                if input[0] <= 31.994999885559082 {
                    var39 = []float64{0.0, 1.0, 0.0, 0.0}
                } else {
                    if input[1] <= 72.1500015258789 {
                        var39 = []float64{0.0, 1.0, 0.0, 0.0}
                    } else {
                        var39 = []float64{0.0, 0.0, 1.0, 0.0}
                    }
                }
            } else {
                if input[1] <= 150.37999725341797 {
                    var39 = []float64{0.0, 0.0, 1.0, 0.0}
                } else {
                    var39 = []float64{0.0, 0.0, 0.0, 1.0}
                }
            }
        }
    }
    var var40 []float64
    if input[1] <= 25.019999504089355 {
        var40 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[0] <= 32.010000228881836 {
            if input[0] <= 22.994999885559082 {
                if input[0] <= 20.015000343322754 {
                    var40 = []float64{0.0, 0.0, 0.0, 1.0}
                } else {
                    var40 = []float64{0.0, 0.0, 1.0, 0.0}
                }
            } else {
                var40 = []float64{0.0, 1.0, 0.0, 0.0}
            }
        } else {
            if input[0] <= 33.9950008392334 {
                var40 = []float64{0.0, 0.0, 1.0, 0.0}
            } else {
                if input[1] <= 138.0449981689453 {
                    var40 = []float64{0.0, 0.0, 1.0, 0.0}
                } else {
                    var40 = []float64{0.0, 0.0, 0.0, 1.0}
                }
            }
        }
    }
    var var41 []float64
    if input[0] <= 22.994999885559082 {
        if input[0] <= 20.015000343322754 {
            var41 = []float64{0.0, 0.0, 0.0, 1.0}
        } else {
            var41 = []float64{0.0, 0.0, 1.0, 0.0}
        }
    } else {
        if input[0] <= 32.010000228881836 {
            if input[0] <= 29.979999542236328 {
                if input[1] <= 24.994999885559082 {
                    var41 = []float64{1.0, 0.0, 0.0, 0.0}
                } else {
                    var41 = []float64{0.0, 1.0, 0.0, 0.0}
                }
            } else {
                var41 = []float64{0.0, 1.0, 0.0, 0.0}
            }
        } else {
            if input[1] <= 150.31500244140625 {
                var41 = []float64{0.0, 0.0, 1.0, 0.0}
            } else {
                var41 = []float64{0.0, 0.0, 0.0, 1.0}
            }
        }
    }
    var var42 []float64
    if input[1] <= 24.984999656677246 {
        var42 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[0] <= 34.0049991607666 {
            if input[1] <= 50.09499931335449 {
                var42 = []float64{0.0, 1.0, 0.0, 0.0}
            } else {
                if input[0] <= 20.074999809265137 {
                    var42 = []float64{0.0, 0.0, 0.0, 1.0}
                } else {
                    var42 = []float64{0.0, 0.0, 1.0, 0.0}
                }
            }
        } else {
            var42 = []float64{0.0, 0.0, 0.0, 1.0}
        }
    }
    var var43 []float64
    if input[0] <= 32.010000228881836 {
        if input[1] <= 24.984999656677246 {
            var43 = []float64{1.0, 0.0, 0.0, 0.0}
        } else {
            if input[0] <= 22.994999885559082 {
                if input[1] <= 150.95999908447266 {
                    var43 = []float64{0.0, 0.0, 1.0, 0.0}
                } else {
                    var43 = []float64{0.0, 0.0, 0.0, 1.0}
                }
            } else {
                if input[1] <= 81.6500015258789 {
                    var43 = []float64{0.0, 1.0, 0.0, 0.0}
                } else {
                    var43 = []float64{0.0, 0.0, 1.0, 0.0}
                }
            }
        }
    } else {
        if input[0] <= 34.0 {
            var43 = []float64{0.0, 0.0, 1.0, 0.0}
        } else {
            var43 = []float64{0.0, 0.0, 0.0, 1.0}
        }
    }
    var var44 []float64
    if input[0] <= 22.994999885559082 {
        if input[1] <= 151.59000396728516 {
            var44 = []float64{0.0, 0.0, 1.0, 0.0}
        } else {
            var44 = []float64{0.0, 0.0, 0.0, 1.0}
        }
    } else {
        if input[0] <= 32.01499938964844 {
            if input[1] <= 24.994999885559082 {
                var44 = []float64{1.0, 0.0, 0.0, 0.0}
            } else {
                var44 = []float64{0.0, 1.0, 0.0, 0.0}
            }
        } else {
            if input[1] <= 152.6699981689453 {
                var44 = []float64{0.0, 0.0, 1.0, 0.0}
            } else {
                var44 = []float64{0.0, 0.0, 0.0, 1.0}
            }
        }
    }
    var var45 []float64
    if input[1] <= 24.97000026702881 {
        var45 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[1] <= 50.119998931884766 {
            var45 = []float64{0.0, 1.0, 0.0, 0.0}
        } else {
            if input[0] <= 33.9950008392334 {
                if input[1] <= 152.58499908447266 {
                    var45 = []float64{0.0, 0.0, 1.0, 0.0}
                } else {
                    var45 = []float64{0.0, 0.0, 0.0, 1.0}
                }
            } else {
                var45 = []float64{0.0, 0.0, 0.0, 1.0}
            }
        }
    }
    var var46 []float64
    if input[1] <= 24.984999656677246 {
        var46 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[0] <= 34.0049991607666 {
            if input[0] <= 20.104999542236328 {
                var46 = []float64{0.0, 0.0, 0.0, 1.0}
            } else {
                if input[0] <= 32.010000228881836 {
                    if input[0] <= 22.994999885559082 {
                        var46 = []float64{0.0, 0.0, 1.0, 0.0}
                    } else {
                        if input[1] <= 81.6500015258789 {
                            var46 = []float64{0.0, 1.0, 0.0, 0.0}
                        } else {
                            var46 = []float64{0.0, 0.0, 1.0, 0.0}
                        }
                    }
                } else {
                    var46 = []float64{0.0, 0.0, 1.0, 0.0}
                }
            }
        } else {
            var46 = []float64{0.0, 0.0, 0.0, 1.0}
        }
    }
    var var47 []float64
    if input[0] <= 22.994999885559082 {
        if input[0] <= 20.015000343322754 {
            var47 = []float64{0.0, 0.0, 0.0, 1.0}
        } else {
            var47 = []float64{0.0, 0.0, 1.0, 0.0}
        }
    } else {
        if input[0] <= 32.010000228881836 {
            if input[1] <= 24.984999656677246 {
                var47 = []float64{1.0, 0.0, 0.0, 0.0}
            } else {
                if input[1] <= 81.67500114440918 {
                    var47 = []float64{0.0, 1.0, 0.0, 0.0}
                } else {
                    var47 = []float64{0.0, 0.0, 1.0, 0.0}
                }
            }
        } else {
            if input[0] <= 33.96999931335449 {
                var47 = []float64{0.0, 0.0, 1.0, 0.0}
            } else {
                if input[1] <= 136.31999969482422 {
                    var47 = []float64{0.0, 0.0, 1.0, 0.0}
                } else {
                    var47 = []float64{0.0, 0.0, 0.0, 1.0}
                }
            }
        }
    }
    var var48 []float64
    if input[1] <= 24.984999656677246 {
        var48 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[0] <= 33.9950008392334 {
            if input[0] <= 22.994999885559082 {
                if input[0] <= 20.015000343322754 {
                    var48 = []float64{0.0, 0.0, 0.0, 1.0}
                } else {
                    var48 = []float64{0.0, 0.0, 1.0, 0.0}
                }
            } else {
                if input[1] <= 50.119998931884766 {
                    var48 = []float64{0.0, 1.0, 0.0, 0.0}
                } else {
                    var48 = []float64{0.0, 0.0, 1.0, 0.0}
                }
            }
        } else {
            var48 = []float64{0.0, 0.0, 0.0, 1.0}
        }
    }
    var var49 []float64
    if input[1] <= 24.994999885559082 {
        var49 = []float64{1.0, 0.0, 0.0, 0.0}
    } else {
        if input[0] <= 22.994999885559082 {
            if input[0] <= 20.015000343322754 {
                var49 = []float64{0.0, 0.0, 0.0, 1.0}
            } else {
                var49 = []float64{0.0, 0.0, 1.0, 0.0}
            }
        } else {
            if input[1] <= 50.18499946594238 {
                var49 = []float64{0.0, 1.0, 0.0, 0.0}
            } else {
                if input[0] <= 34.01499938964844 {
                    if input[1] <= 358.2200012207031 {
                        var49 = []float64{0.0, 0.0, 1.0, 0.0}
                    } else {
                        var49 = []float64{0.0, 0.0, 0.0, 1.0}
                    }
                } else {
                    var49 = []float64{0.0, 0.0, 0.0, 1.0}
                }
            }
        }
    }
    return mulVectorNumberOverall(addVectorsOverall(addVectorsOverall(addVectorsOverall(addVectorsOverall(addVectorsOverall(addVectorsOverall(addVectorsOverall(addVectorsOverall(addVectorsOverall(addVectorsOverall(addVectorsOverall(addVectorsOverall(addVectorsOverall(addVectorsOverall(addVectorsOverall(addVectorsOverall(addVectorsOverall(addVectorsOverall(addVectorsOverall(addVectorsOverall(addVectorsOverall(addVectorsOverall(addVectorsOverall(addVectorsOverall(addVectorsOverall(addVectorsOverall(addVectorsOverall(addVectorsOverall(addVectorsOverall(addVectorsOverall(addVectorsOverall(addVectorsOverall(addVectorsOverall(addVectorsOverall(addVectorsOverall(addVectorsOverall(addVectorsOverall(addVectorsOverall(addVectorsOverall(addVectorsOverall(addVectorsOverall(addVectorsOverall(addVectorsOverall(addVectorsOverall(addVectorsOverall(addVectorsOverall(addVectorsOverall(addVectorsOverall(addVectorsOverall(var0, var1), var2), var3), var4), var5), var6), var7), var8), var9), var10), var11), var12), var13), var14), var15), var16), var17), var18), var19), var20), var21), var22), var23), var24), var25), var26), var27), var28), var29), var30), var31), var32), var33), var34), var35), var36), var37), var38), var39), var40), var41), var42), var43), var44), var45), var46), var47), var48), var49), 0.02)
}
func addVectorsOverall(v1, v2 []float64) []float64 {
    result := make([]float64, len(v1))
    for i := 0; i < len(v1); i++ {
        result[i] = v1[i] + v2[i]
    }
    return result
}
func mulVectorNumberOverall(v1 []float64, num float64) []float64 {
    result := make([]float64, len(v1))
    for i := 0; i < len(v1); i++ {
        result[i] = v1[i] * num
    }
    return result
}
