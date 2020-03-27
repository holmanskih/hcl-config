  Standard HCL2 **hclsimple.DecodeFile** function
* provides parsing of one file or many files but with the same structure that 
  can be overwritten by the previous .hcl file!
* if the structure of e.g. two file is different -> throws an error
* all golang structures need to have the label/block meta info
* all golang structures need to have additional fields if the configuration
  is of type block/block with label
  
  
  